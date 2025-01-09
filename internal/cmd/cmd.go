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

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
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
	Short: "Create a backup of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		dbService := database.NewDatabaseService()
		dbService.Backup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(backupCmd)
}
