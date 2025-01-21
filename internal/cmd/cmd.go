package cmd

import (
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/database"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "go-vault",
	Short: "CLI tool for db backups and restores",
	Long:  "Manage your database backups and restores with ease.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup config of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		configService := config.NewConfigService()
		configService.SetupConfig()
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups of your database.",
}

var createBackupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		dbService := database.NewDatabaseService()
		dbService.CreateBackup()
	},
}

var listBackupCmd = &cobra.Command{
	Use: "list",
	Aliases: []string{
		"ls",
	},
	Short: "List all backups of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		dbService := database.NewDatabaseService()
		dbService.ListBackups()
	},
}

var restoreBackupCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		dbService := database.NewDatabaseService()
		dbService.RestoreBackup("1737487435_postgres_backup.sql")
	},
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
