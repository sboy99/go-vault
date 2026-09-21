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

func rootCmdHandler(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func setupCmdHandler(cmd *cobra.Command, args []string) {
	config.LoadOptional()
	if err := setup.NewConfigService().SetupConfig(); err != nil {
		logger.Error("%v", err)
	}
}

func backupCmdHandler(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func createBackupCmdHandler(cmd *cobra.Command, args []string) {
	svc, err := newBackupService()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	b, err := svc.CreateBackup(context.Background())
	if err != nil {
		logger.Error("backup failed: %v", err)
		return
	}
	logger.Info("created backup %s (%s)", b.BackupId, b.Name)
}

func listBackupCmdHandler(cmd *cobra.Command, args []string) {
	svc, err := newBackupService()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	list, err := svc.ListBackups(50, 0)
	if err != nil {
		logger.Error("list failed: %v", err)
		return
	}
	headers, err := utils.GetStructFields(meta.BackupMeta{})
	if err != nil {
		logger.Error("%v", err)
		return
	}
	rows := make([]interface{}, len(list))
	for i, v := range list {
		rows[i] = *v
	}
	if err := ui.RenderTable(headers, rows); err != nil {
		logger.Error("%v", err)
	}
}

func restoreBackupCmdHandler(cmd *cobra.Command, args []string) {
	svc, err := newBackupService()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	if err := svc.RestoreBackup(context.Background(), args[0]); err != nil {
		logger.Error("restore failed: %v", err)
		return
	}
	logger.Info("restore complete")
}

func serveCmdHandler(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Fatal("%v", err)
	}
	svc := backup.NewService(cfg, store)
	runner := job.NewRunner()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := svc.Reconcile(ctx); err != nil {
		logger.Warn("reconcile: %v", err)
	}
	cancel()

	sched, err := scheduler.New(cfg, svc, runner)
	if err != nil {
		logger.Fatal("%v", err)
	}
	if err := sched.Start(); err != nil {
		logger.Fatal("%v", err)
	}

	metrics.SetReady(true)
	server := api.NewServer(cfg, svc, runner)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

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

	metrics.SetReady(false)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Runtime.ShutdownTimeout)
	defer shutdownCancel()
	sched.Stop(shutdownCtx)
	runner.Cancel()
	_ = server.Shutdown(shutdownCtx)
}

func newBackupService() (*backup.Service, error) {
	cfg := config.GetConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	return backup.NewService(cfg, store), nil
}
