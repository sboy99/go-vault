package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/pkg/logger"
)

// BackupConfig is the narrow config surface BackupService needs.
type BackupConfig struct {
	DBName         string
	DatabaseType   domain.DatabaseType
	StorageType    domain.StorageType
	Conn           domain.ConnParams
	TempDir        string
	DumpTimeout    time.Duration
	RestoreTimeout time.Duration
	RestoreJobs    int
	Retention      domain.RetentionPolicy
	Location       *time.Location
}

// BackupService orchestrates dump, store, verify, restore, prune, and reconcile.
type BackupService struct {
	cfg      BackupConfig
	backups  domain.BackupRepository
	store    domain.ArtifactStore
	engine   domain.DumpEngine
	notifier domain.Notifier
	metrics  domain.MetricsRecorder
}

func NewBackupService(
	cfg BackupConfig,
	backups domain.BackupRepository,
	store domain.ArtifactStore,
	eng domain.DumpEngine,
	notifier domain.Notifier,
	rec domain.MetricsRecorder,
) *BackupService {
	if cfg.Location == nil {
		cfg.Location = time.UTC
	}
	return &BackupService{
		cfg:      cfg,
		backups:  backups,
		store:    store,
		engine:   eng,
		notifier: notifier,
		metrics:  rec,
	}
}

func (s *BackupService) CreateBackup(ctx context.Context) (*domain.Backup, error) {
	started := time.Now().UTC()
	backup, err := s.startBackup(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.DumpTimeout)
	defer cancel()

	info, dumpRes, err := s.dumpAndSave(ctx, backup.StorageKey)
	if err != nil {
		return s.failBackup(ctx, backup, "dump/save", err)
	}
	if err := s.verifyDump(ctx, backup, dumpRes); err != nil {
		_ = s.store.Delete(ctx, backup.StorageKey)
		return s.failBackup(ctx, backup, "verify", err)
	}
	return s.markSuccess(ctx, backup, info, dumpRes, started)
}

func (s *BackupService) BackupAndPrune(ctx context.Context) (*domain.Backup, error) {
	b, err := s.CreateBackup(ctx)
	if err != nil {
		return b, err
	}
	if _, err := s.Prune(ctx); err != nil {
		return b, err
	}
	return b, nil
}

func (s *BackupService) RestoreBackup(ctx context.Context, backupID, confirm string) error {
	if confirm != s.cfg.DBName {
		return domain.ErrRestoreConfirmMismatch
	}

	backup, err := s.resolveBackup(ctx, backupID)
	if err != nil {
		return err
	}
	if backup.Status != domain.StatusSuccess {
		return fmt.Errorf("%w: status %s", domain.ErrBackupNotRestorable, backup.Status)
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.RestoreTimeout)
	defer cancel()

	if err := os.MkdirAll(s.cfg.TempDir, 0o755); err != nil {
		return err
	}
	tmpPath := filepath.Join(s.cfg.TempDir, backup.Name+".restore")
	if err := s.spoolToFile(ctx, backup.StorageKey, tmpPath); err != nil {
		s.metrics.RestoreFinished("failed")
		return err
	}
	defer func() { _ = os.Remove(tmpPath) }()

	if err := s.engine.Restore(ctx, s.cfg.Conn, tmpPath, s.cfg.RestoreJobs); err != nil {
		s.metrics.RestoreFinished("failed")
		s.notifier.NotifyFailure(ctx, "restore_failed", err.Error())
		return err
	}
	s.metrics.RestoreFinished("success")
	logger.Info("restore successful backup_id=%s", backup.BackupId)
	return nil
}

func (s *BackupService) ListBackups(ctx context.Context, limit, offset int) ([]*domain.Backup, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.backups.List(ctx, limit, offset)
}

func (s *BackupService) GetBackup(ctx context.Context, id string) (*domain.Backup, error) {
	return s.backups.FindByID(ctx, id)
}

func (s *BackupService) OpenBackup(ctx context.Context, id string) (io.ReadCloser, *domain.Backup, error) {
	b, err := s.backups.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	rc, err := s.store.Open(ctx, b.StorageKey)
	if err != nil {
		return nil, nil, err
	}
	return rc, b, nil
}

func (s *BackupService) Prune(ctx context.Context) (int, error) {
	all, err := s.backups.FindAll(ctx)
	if err != nil {
		return 0, err
	}
	flat := make([]domain.Backup, 0, len(all))
	for _, b := range all {
		flat = append(flat, *b)
	}

	_, prune := domain.Select(flat, s.cfg.Retention, time.Now().UTC(), s.cfg.Location)

	pruned := 0
	for _, b := range prune {
		if err := s.store.Delete(ctx, b.StorageKey); err != nil {
			logger.Warn("prune storage delete failed for %s: %v", b.StorageKey, err)
			continue
		}
		if err := s.backups.Delete(ctx, b.BackupId); err != nil {
			logger.Warn("prune meta delete failed for %s: %v", b.BackupId, err)
			continue
		}
		pruned++
		s.metrics.BackupPruned()
	}
	logger.Info("pruned %d backups", pruned)
	return pruned, nil
}

func (s *BackupService) Reconcile(ctx context.Context) error {
	objects, err := s.store.List(ctx, "")
	if err != nil {
		return err
	}
	existing, err := s.backups.FindAll(ctx)
	if err != nil {
		return err
	}
	known := map[string]struct{}{}
	for _, b := range existing {
		known[b.StorageKey] = struct{}{}
		if b.Name != "" {
			known[b.Name] = struct{}{}
		}
	}
	for _, obj := range objects {
		key := obj.Key
		base := filepath.Base(key)
		if _, ok := known[key]; ok {
			continue
		}
		if _, ok := known[base]; ok {
			continue
		}
		if !strings.HasSuffix(base, ".dump") && !strings.HasSuffix(base, ".sql") {
			continue
		}
		m := domain.NewBackup(base, s.cfg.DatabaseType, s.cfg.StorageType)
		m.StorageKey = key
		m.Status = domain.StatusSuccess
		m.SizeBytes = obj.Size
		m.CreatedAt = obj.LastModified
		m.StartedAt = obj.LastModified
		finished := obj.LastModified
		m.FinishedAt = &finished
		if err := s.backups.Save(ctx, m); err != nil {
			logger.Warn("reconcile save failed for %s: %v", key, err)
			continue
		}
		logger.Info("reconciled missing metadata for %s", key)
	}
	return nil
}

func (s *BackupService) Ping(ctx context.Context) error {
	_, _, err := s.engine.ServerVersion(ctx, s.cfg.Conn)
	return err
}

func (s *BackupService) startBackup(ctx context.Context) (*domain.Backup, error) {
	filename := buildFileName(s.cfg.DBName)
	backup := domain.NewBackup(filename, s.cfg.DatabaseType, s.cfg.StorageType)
	backup.StorageKey = filename
	if err := s.backups.Save(ctx, backup); err != nil {
		return nil, fmt.Errorf("save running meta: %w", err)
	}
	return backup, nil
}

func (s *BackupService) dumpAndSave(ctx context.Context, filename string) (domain.ObjectInfo, *domain.DumpResult, error) {
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)
	var dumpRes *domain.DumpResult

	go func() {
		defer func() { _ = pw.Close() }()
		res, err := s.engine.DumpTo(ctx, s.cfg.Conn, pw)
		dumpRes = res
		errCh <- err
	}()

	info, saveErr := s.store.Save(ctx, filename, pr)
	dumpErr := <-errCh
	if dumpErr != nil {
		return domain.ObjectInfo{}, dumpRes, fmt.Errorf("dump: %w", dumpErr)
	}
	if saveErr != nil {
		return domain.ObjectInfo{}, dumpRes, fmt.Errorf("save: %w", saveErr)
	}
	return info, dumpRes, nil
}

func (s *BackupService) verifyDump(ctx context.Context, backup *domain.Backup, dumpRes *domain.DumpResult) error {
	if err := os.MkdirAll(s.cfg.TempDir, 0o755); err != nil {
		return err
	}
	tmpPath := filepath.Join(s.cfg.TempDir, backup.StorageKey+".verify")
	if err := s.spoolToFile(ctx, backup.StorageKey, tmpPath); err != nil {
		return fmt.Errorf("spool for verify: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	major := 0
	if dumpRes != nil {
		major = dumpRes.ServerMajor
	}
	return s.engine.Verify(ctx, tmpPath, major)
}

func (s *BackupService) failBackup(ctx context.Context, backup *domain.Backup, stage string, err error) (*domain.Backup, error) {
	failed := backup.MarkFailed(err.Error())
	_ = s.backups.Save(ctx, &failed)
	s.metrics.BackupFailed()
	s.notifier.NotifyFailure(ctx, "backup_failed", err.Error())
	return &failed, fmt.Errorf("%s: %w", stage, err)
}

func (s *BackupService) markSuccess(ctx context.Context, backup *domain.Backup, info domain.ObjectInfo, dumpRes *domain.DumpResult, started time.Time) (*domain.Backup, error) {
	pgVersion := ""
	if dumpRes != nil {
		pgVersion = dumpRes.ServerVersion
	}
	success := backup.MarkSuccess(info.Size, info.SHA256, pgVersion, true)
	if err := s.backups.Save(ctx, &success); err != nil {
		return &success, err
	}
	s.metrics.BackupSucceeded(started, info.Size)
	logger.Info("backup successful id=%s size=%d", success.BackupId, info.Size)
	return &success, nil
}

func (s *BackupService) resolveBackup(ctx context.Context, backupID string) (*domain.Backup, error) {
	backup, err := s.backups.FindByID(ctx, backupID)
	if err == nil {
		return backup, nil
	}
	all, listErr := s.backups.FindAll(ctx)
	if listErr != nil {
		return nil, err
	}
	for _, b := range all {
		if b.Name == backupID || b.StorageKey == backupID {
			return b, nil
		}
	}
	return nil, domain.ErrBackupNotFound
}

func (s *BackupService) spoolToFile(ctx context.Context, key, dest string) error {
	rc, err := s.store.Open(ctx, key)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = io.Copy(f, rc)
	return err
}

func buildFileName(dbName string) string {
	return fmt.Sprintf("%d_%s_backup.dump", time.Now().Unix(), dbName)
}
