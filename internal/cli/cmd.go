package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-vault",
	Short: "Production PostgreSQL backup service powered by pg_dump",
	Long:  "Schedule, store, list, and restore PostgreSQL backups using the official pg_dump/pg_restore binaries.",
	RunE:  rootCmdHandler,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup of database and storage config.",
	RunE:  setupCmdHandler,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the backup scheduler and HTTP API.",
	RunE:  serveCmdHandler,
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups of your database.",
	RunE:  backupCmdHandler,
}

var createBackupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of your database.",
	RunE:  createBackupCmdHandler,
}

var listBackupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups of your database.",
	RunE:  listBackupCmdHandler,
}

var restoreBackupCmd = &cobra.Command{
	Use:               "restore [backup_id_or_name]",
	Short:             "Restore a backup of your database.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: restoreBackupValidArgs,
	RunE:              restoreBackupCmdHandler,
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(listBackupCmd)
	backupCmd.AddCommand(restoreBackupCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
