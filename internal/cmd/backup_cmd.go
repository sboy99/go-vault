package cmd

import (
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/sboy99/go-vault/pkg/logger"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		dbType, err := ui.DisplaySelectDBPrompt()
		if err != nil {
			logger.Error("Error selecting DB type: %v", err)
		}
		dbConfig, err := ui.DisplayDBConfigPrompt()
		if err != nil {
			logger.Error("Error getting DB config: %v", err)
		}

		logger.Info("Selected DB: %v", dbType)
		logger.Info("DB Config: %v", dbConfig)

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
