package cli

import (
	"context"
	"fmt"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/sboy99/go-vault/internal/utils"
	"github.com/sboy99/go-vault/internal/version"
	"github.com/sboy99/go-vault/pkg/logger"
	"github.com/spf13/cobra"
)

// Deps holds services the CLI transport needs.
type Deps struct {
	Backup *app.BackupService
	Jobs   *app.JobService
	Setup  *app.SetupService
}

// NewRootCommand builds the Cobra command tree with injected services.
// The serve command is intentionally omitted — that lives in cmd/server.
func NewRootCommand(deps Deps) *cobra.Command {
	root := &cobra.Command{
		Use:     "go-vault",
		Short:   "Production PostgreSQL backup service powered by pg_dump",
		Long:    "Schedule, store, list, and restore PostgreSQL backups using official pg_dump (plain SQL + gzip) and psql.",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(newSetupCommand(deps.Setup))
	root.AddCommand(newBackupCommand(deps))
	return root
}

func newBackupCommand(deps Deps) *cobra.Command {
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage backups of your database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a backup of your database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := deps.Backup.BackupAndPrune(context.Background())
			if err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
			logger.Info("created backup %s (%s)", b.BackupId, b.Name)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all backups of your database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := deps.Backup.ListBackups(context.Background(), 50, 0)
			if err != nil {
				return fmt.Errorf("list failed: %w", err)
			}
			headers, err := utils.GetStructFields(domain.Backup{})
			if err != nil {
				return err
			}
			rows := make([]interface{}, len(list))
			for i, v := range list {
				rows[i] = *v
			}
			return ui.RenderTable(headers, rows)
		},
	}

	restoreCmd := &cobra.Command{
		Use:               "restore [backup_id_or_name]",
		Short:             "Restore a backup of your database.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: restoreValidArgs(deps.Backup),
		RunE: func(cmd *cobra.Command, args []string) error {
			confirm, err := ui.DisplayRestoreConfirmPrompt()
			if err != nil {
				return fmt.Errorf("confirm database name: %w", err)
			}
			if err := deps.Backup.RestoreBackup(context.Background(), args[0], confirm); err != nil {
				return fmt.Errorf("restore failed: %w", err)
			}
			logger.Info("restore complete")
			return nil
		},
	}

	backupCmd.AddCommand(createCmd, listCmd, restoreCmd)
	return backupCmd
}
