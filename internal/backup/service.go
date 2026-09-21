package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sboy99/go-vault/internal/alert"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/engine"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/internal/retention"
	"github.com/sboy99/go-vault/internal/utils"
	"github.com/sboy99/go-vault/pkg/logger"
)

// Service orchestrates dump, store, verify, restore, and prune.
type Service struct {
	cfg      *config.Config
	store    domain.ArtifactStore
	engine   domain.DumpEngine
	notifier domain.Notifier
	metrics  domain.MetricsRecorder
}

func NewService(cfg *config.Config, store domain.ArtifactStore) *Service {
	return &Service{
		cfg:      cfg,
		store:    store,
		engine:   engine.NewPostgresEngine(),
		notifier: alert.NewWebhookNotifier(cfg.Runtime.AlertWebhook),
		metrics:  metrics.NewPrometheusRecorder(),
	}
}

func (s *Service) CreateBackup(ctx context.Context) (*meta.BackupMeta, error) {
	started := time.Now().UTC()
	backupMeta, err := s.startMeta()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Runtime.DumpTimeout)
	defer cancel()

	info, dumpRes, err := s.dumpAndSave(ctx, backupMeta.StorageKey)
	if err != nil {
		return s.failBackup(backupMeta, "dump/save", err)
	}
	if err := s.verifyDump(ctx, backupMeta, dumpRes); err != nil {
		_ = s.store.Delete(ctx, backupMeta.StorageKey)
		return s.failBackup(backupMeta, "verify", err)
	}
	return s.markSuccess(backupMeta, info, dumpRes, started)
}

func (s *Service) BackupAndPrune(ctx context.Context) (*meta.BackupMeta, error) {
	b, err := s.CreateBackup(ctx)
	if err != nil {
		return b, err
	}
	if _, err := s.Prune(ctx); err != nil {
		return b, err
	}
	return b, nil
}

func (s *Service) RestoreBackup(ctx context.Context, backupID string) error {
	backupMeta, err := s.resolveBackup(backupID)
	if err != nil {
		return err
	}
	if backupMeta.Status != meta.StatusSuccess {
		return fmt.Errorf("cannot restore backup in status %s", backupMeta.Status)
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.Runtime.RestoreTimeout)
	defer cancel()

	if err := os.MkdirAll(s.cfg.Runtime.TempDir, 0o755); err != nil {
		return err
	}
	tmpPath := filepath.Join(s.cfg.Runtime.TempDir, backupMeta.Name+".restore")
	if err := s.spoolToFile(ctx, backupMeta.StorageKey, tmpPath); err != nil {
		s.metrics.RestoreFinished("failed")
		return err
	}
	defer func() { _ = os.Remove(tmpPath) }()

	if err := s.engine.Restore(ctx, s.connParams(), tmpPath, s.cfg.Runtime.RestoreJobs); err != nil {
		s.metrics.RestoreFinished("failed")
		s.notifier.NotifyFailure(ctx, "restore_failed", err.Error())
		return err
	}
	s.metrics.RestoreFinished("success")
	logger.Info("restore successful backup_id=%s", backupMeta.BackupId)
	return nil
}

func (s *Service) ListBackups(limit, offset int) ([]*meta.BackupMeta, error) {
	if limit <= 0 {
		limit = 50
	}
	return meta.ListBackupMeta(limit, offset)
}

func (s *Service) GetBackup(id string) (*meta.BackupMeta, error) {
	return meta.GetBackupMeta(id)
}

func (s *Service) OpenBackup(ctx context.Context, id string) (io.ReadCloser, *meta.BackupMeta, error) {
	b, err := meta.GetBackupMeta(id)
	if err != nil {
		return nil, nil, err
	}
	rc, err := s.store.Open(ctx, b.StorageKey)
	if err != nil {
		return nil, nil, err
	}
	return rc, b, nil
}

func (s *Service) Prune(ctx context.Context) (int, error) {
	all, err := meta.ListAllBackupMeta()
	if err != nil {
		return 0, err
	}
	flat := make([]meta.BackupMeta, 0, len(all))
	for _, b := range all {
		flat = append(flat, *b)
	}

	loc, err := time.LoadLocation(s.cfg.Schedule.Timezone)
	if err != nil {
		loc = time.UTC
	}
	policy := retention.Policy{
		Daily:   s.cfg.Retention.Daily,
		Weekly:  s.cfg.Retention.Weekly,
		Monthly: s.cfg.Retention.Monthly,
	}
	_, prune := retention.Select(flat, policy, time.Now().UTC(), loc)

	pruned := 0
	for _, b := range prune {
		if err := s.store.Delete(ctx, b.StorageKey); err != nil {
			logger.Warn("prune storage delete failed for %s: %v", b.StorageKey, err)
			continue
		}
		if err := meta.DeleteBackupMeta(b.BackupId); err != nil {
			logger.Warn("prune meta delete failed for %s: %v", b.BackupId, err)
			continue
		}
		pruned++
		s.metrics.BackupPruned()
	}
	logger.Info("pruned %d backups", pruned)
	return pruned, nil
}

func (s *Service) Reconcile(ctx context.Context) error {
	return meta.Reconcile(ctx, func(ctx context.Context, prefix string) ([]meta.ObjectRef, error) {
		objs, err := s.store.List(ctx, prefix)
		if err != nil {
			return nil, err
		}
		out := make([]meta.ObjectRef, 0, len(objs))
		for _, o := range objs {
			out = append(out, meta.ObjectRef{
				Key:          o.Key,
				Size:         o.Size,
				LastModified: o.LastModified,
			})
		}
		return out, nil
	}, s.cfg.DB.Type, s.cfg.Storage.Type)
}

func (s *Service) Ping(ctx context.Context) error {
	_, _, err := s.engine.ServerVersion(ctx, s.connParams())
	return err
}

func (s *Service) startMeta() (*meta.BackupMeta, error) {
	filename := buildFileName(s.cfg.DB.Name)
	backupMeta := meta.NewBackupMeta(filename, s.cfg.DB.Type, s.cfg.Storage.Type)
	backupMeta.StorageKey = filename
	if err := backupMeta.Save(); err != nil {
		return nil, fmt.Errorf("save running meta: %w", err)
	}
	return backupMeta, nil
}

func (s *Service) dumpAndSave(ctx context.Context, filename string) (domain.ObjectInfo, *domain.DumpResult, error) {
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)
	var dumpRes *domain.DumpResult

	go func() {
		defer func() { _ = pw.Close() }()
		res, err := s.engine.DumpTo(ctx, s.connParams(), pw)
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

func (s *Service) verifyDump(ctx context.Context, backupMeta *meta.BackupMeta, dumpRes *domain.DumpResult) error {
	if err := os.MkdirAll(s.cfg.Runtime.TempDir, 0o755); err != nil {
		return err
	}
	tmpPath := filepath.Join(s.cfg.Runtime.TempDir, backupMeta.StorageKey+".verify")
	if err := s.spoolToFile(ctx, backupMeta.StorageKey, tmpPath); err != nil {
		return fmt.Errorf("spool for verify: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	major := 0
	if dumpRes != nil {
		major = dumpRes.ServerMajor
	}
	return s.engine.Verify(ctx, tmpPath, major)
}

func (s *Service) failBackup(backupMeta *meta.BackupMeta, stage string, err error) (*meta.BackupMeta, error) {
	_ = backupMeta.MarkFailed(err.Error())
	s.metrics.BackupFailed()
	s.notifier.NotifyFailure(context.Background(), "backup_failed", err.Error())
	return backupMeta, fmt.Errorf("%s: %w", stage, err)
}

func (s *Service) markSuccess(backupMeta *meta.BackupMeta, info domain.ObjectInfo, dumpRes *domain.DumpResult, started time.Time) (*meta.BackupMeta, error) {
	pgVersion := ""
	if dumpRes != nil {
		pgVersion = dumpRes.ServerVersion
	}
	if err := backupMeta.MarkSuccess(info.Size, info.SHA256, pgVersion, true); err != nil {
		return backupMeta, err
	}
	s.metrics.BackupSucceeded(started, info.Size)
	logger.Info("backup successful id=%s size=%d", backupMeta.BackupId, info.Size)
	return backupMeta, nil
}

func (s *Service) resolveBackup(backupID string) (*meta.BackupMeta, error) {
	backupMeta, err := meta.GetBackupMeta(backupID)
	if err == nil {
		return backupMeta, nil
	}
	all, listErr := meta.ListAllBackupMeta()
	if listErr != nil {
		return nil, err
	}
	for _, b := range all {
		if b.Name == backupID || b.StorageKey == backupID {
			return b, nil
		}
	}
	return nil, fmt.Errorf("backup not found: %s", backupID)
}

func (s *Service) connParams() domain.ConnParams {
	return domain.ConnParams{
		Host:     s.cfg.DB.Host,
		Port:     s.cfg.DB.Port,
		Name:     s.cfg.DB.Name,
		Username: s.cfg.DB.Username,
		Password: s.cfg.DB.Password,
		SSLMode:  s.cfg.DB.SSLMode,
	}
}

func (s *Service) spoolToFile(ctx context.Context, key, dest string) error {
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
	return fmt.Sprintf("%s_%s_backup.dump", utils.GetUnixTimeStamp(), dbName)
}
