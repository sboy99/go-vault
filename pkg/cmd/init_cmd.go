package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-nester/pkg/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Nestjs application",
	Long:  "Initialize a new Nestjs application in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.NewConfig()
		// Get the architecture type
		architecture, err := getArchitecture()
		if err != nil {
			panic(err)
		}
		config.Arcitecture = architecture
		// Create the config file
		err = config.Save()
		if err != nil {
			panic(err)
		}
	},
}

func getArchitecture() (config.Arcitecture, error) {
	arcitectures := []config.Arcitecture{config.HEXAGONAL}
	prompt := promptui.Select{
		Label: "Select the architecture type",
		Items: arcitectures,
	}
	_, architecture, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.Arcitecture(architecture), nil
}
