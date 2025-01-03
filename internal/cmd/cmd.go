package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "go-vault",
	Short: "CLI tool for db backups and restores",
	Long:  "Manage your database backups and restores with ease.",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
