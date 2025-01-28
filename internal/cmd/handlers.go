package cmd

import (
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/database"
	"github.com/spf13/cobra"
)

func rootCmdHandler(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func setupCmdHandler(cmd *cobra.Command, args []string) {
	configService := config.NewConfigService()
	configService.SetupConfig()
}

func backupCmdHandler(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func createBackupCmdHandler(cmd *cobra.Command, args []string) {
	dbService := database.NewDatabaseService()
	dbService.CreateBackup()
}

func listBackupCmdHandler(cmd *cobra.Command, args []string) {
	dbService := database.NewDatabaseService()
	dbService.ListBackups()
}

func restoreBackupCmdHandler(cmd *cobra.Command, args []string) {
	backupName := args[0]
	dbService := database.NewDatabaseService()
	dbService.RestoreBackup(backupName)
}
