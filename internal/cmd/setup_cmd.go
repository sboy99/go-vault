package cmd

import (
	"github.com/sboy99/go-vault/internal/services"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup config of your database.",
	Run: func(cmd *cobra.Command, args []string) {
		services.SetupConfig()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
