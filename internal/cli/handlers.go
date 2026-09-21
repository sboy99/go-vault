package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sboy99/go-vault/internal/api"
	"github.com/sboy99/go-vault/internal/backup"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/job"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/internal/metrics"
	"github.com/sboy99/go-vault/internal/scheduler"
	"github.com/sboy99/go-vault/internal/setup"
	"github.com/sboy99/go-vault/internal/storage"
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/sboy99/go-vault/internal/utils"
	"github.com/sboy99/go-vault/pkg/logger"
	"github.com/spf13/cobra"
)

func rootCmdHandler(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func setupCmdHandler(cmd *cobra.Command, args []string) error {
	config.LoadOptional()
	return setup.NewConfigService().SetupConfig()
}

func backupCmdHandler(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func createBackupCmdHandler(cmd *cobra.Command, args []string) error {
	svc, err := newBackupService()
	if err != nil {
		return err
	}
	b, err := svc.BackupAndPrune(context.Background())
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}
	logger.Info("created backup %s (%s)", b.BackupId, b.Name)
	return nil
}

func listBackupCmdHandler(cmd *cobra.Command, args []string) error {
	svc, err := newBackupService()
	if err != nil {
		return err
	}
	list, err := svc.ListBackups(50, 0)
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}
	headers, err := utils.GetStructFields(meta.BackupMeta{})
	if err != nil {
		return err
	}
	rows := make([]interface{}, len(list))
	for i, v := range list {
		rows[i] = *v
	}
	return ui.RenderTable(headers, rows)
}

func restoreBackupCmdHandler(cmd *cobra.Command, args []string) error {
	svc, err := newBackupService()
	if err != nil {
		return err
	}
	if err := svc.RestoreBackup(context.Background(), args[0]); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}
	logger.Info("restore complete")
	return nil
}

func serveCmdHandler(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		return err
	}
	svc := backup.NewService(cfg, store)
	runner := job.NewRunner()

	reconcileOnStart(svc)

	sched, err := scheduler.New(cfg, svc, runner)
	if err != nil {
		return err
	}
	if err := sched.Start(); err != nil {
		return err
	}

	metrics.SetReady(true)
	server := api.NewServer(cfg, svc, runner)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	waitForShutdown(errCh)
	return shutdown(cfg, sched, runner, server)
}

func reconcileOnStart(svc *backup.Service) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := svc.Reconcile(ctx); err != nil {
		logger.Warn("reconcile: %v", err)
	}
}

func waitForShutdown(errCh <-chan error) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("received signal %s, shutting down", sig.String())
	case err := <-errCh:
		if err != nil && err.Error() != "http: Server closed" {
			logger.Error("api error: %v", err)
		}
	}
}

func shutdown(cfg *config.Config, sched *scheduler.Scheduler, runner *job.Runner, server *api.Server) error {
	metrics.SetReady(false)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Runtime.ShutdownTimeout)
	defer cancel()
	sched.Stop(ctx)
	runner.Cancel()
	return server.Shutdown(ctx)
}

func newBackupService() (*backup.Service, error) {
	cfg := config.GetConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	return backup.NewService(cfg, store), nil
}
