package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "go-vault",
	Short: "CLI tool for db backups and restores",
	Long:  "Manage your database backups and restores with ease.",
	Run:   rootCmdHandler,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup config of your database.",
	Run:   setupCmdHandler,
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups of your database.",
	Run:   backupCmdHandler,
}

var createBackupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of your database.",
	Run:   createBackupCmdHandler,
}

var listBackupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups of your database.",
	Run:   listBackupCmdHandler,
}

var restoreBackupCmd = &cobra.Command{
	Use:               "restore [backup_name]",
	Short:             "Restore a backup of your database.",
	ValidArgsFunction: restoreBackupValidArgs,
	Run:               restoreBackupCmdHandler,
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(listBackupCmd)
	backupCmd.AddCommand(restoreBackupCmd)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
