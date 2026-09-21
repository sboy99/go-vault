package bootstrap

import (
	"fmt"

	"github.com/sboy99/go-vault/internal/alert"
	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/engine"
	"github.com/sboy99/go-vault/internal/metrics"
	repo "github.com/sboy99/go-vault/internal/repository/boltdb"
	"github.com/sboy99/go-vault/internal/storage"
	boltclient "github.com/sboy99/go-vault/pkg/boltdb"
)

// Container holds wired application services and shared resources.
type Container struct {
	Config  *config.Config
	DB      *boltclient.Client
	Backup  *app.BackupService
	Jobs    *app.JobService
	Setup   *app.SetupService
}

// New builds repositories, adapters, and app services from config.
func New(cfg *config.Config) (*Container, error) {
	db, err := boltclient.New(cfg.Runtime.MetaDBPath)
	if err != nil {
		return nil, fmt.Errorf("boltdb: %w", err)
	}

	backupRepo, err := repo.NewBackupRepository(db)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("backup repository: %w", err)
	}
	jobRepo, err := repo.NewJobRepository(db)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("job repository: %w", err)
	}

	store, err := storage.NewStorage(cfg)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("storage: %w", err)
	}

	backupCfg, err := cfg.ToBackupConfig()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("backup config: %w", err)
	}

	backupSvc := app.NewBackupService(
		backupCfg,
		backupRepo,
		store,
		engine.NewPostgresEngine(),
		alert.NewWebhookNotifier(cfg.Runtime.AlertWebhook),
		metrics.NewPrometheusRecorder(),
	)
	jobSvc := app.NewJobService(jobRepo)
	setupSvc := app.NewSetupService(config.NewSettingsStore())

	return &Container{
		Config: cfg,
		DB:     db,
		Backup: backupSvc,
		Jobs:   jobSvc,
		Setup:  setupSvc,
	}, nil
}

// Close releases shared resources.
func (c *Container) Close() error {
	if c == nil || c.DB == nil {
		return nil
	}
	return c.DB.Close()
}
