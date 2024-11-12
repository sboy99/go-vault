package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "Go Nester",
	Short: "Nestjs CLI tool",
	Long:  "CLI tool for generating files and code snippets for Nestjs applications, written in Go.",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createCmd)
}
