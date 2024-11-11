package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new file or code snippet",
	Long:  "Create a new file or code snippet for your Nestjs application.",
}

var createControllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Create a new controller",
	Long:  "Create a new controller for your Nestjs application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new controller...")
	},
}

var createServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Create a new service",
	Long:  "Create a new service for your Nestjs application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new service...")
	},
}

var createRepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Create a new repository",
	Long:  "Create a new repository for your Nestjs application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new repository...")
	},
}

var createAdapterCmd = &cobra.Command{
	Use:   "adapter",
	Short: "Create a new adapter",
	Long:  "Create a new adapter for your Nestjs application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new adapter...")
	},
}

var createEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Create a new entity",
	Long:  "Create a new entity for your Nestjs application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new entity...")
	},
}

func init() {
	createCmd.AddCommand(createControllerCmd)
	createCmd.AddCommand(createServiceCmd)
	createCmd.AddCommand(createRepositoryCmd)
	createCmd.AddCommand(createAdapterCmd)
	createCmd.AddCommand(createEntityCmd)
}
