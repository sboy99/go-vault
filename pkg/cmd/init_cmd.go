package cmd

import (
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-nester/pkg/config"
	"github.com/sboy99/go-nester/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const _CONFIG_FILE string = "nester.yaml"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Nestjs application",
	Long:  "Initialize a new Nestjs application in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.Config{}
		// Get the architecture type
		architecture, err := getArchitecture()
		if err != nil {
			panic(err)
		}
		config.Arcitecture = architecture
		// Create the config file
		err = createConfigFile(&config)
		if err != nil {
			panic(err)
		}
	},
}

func getArchitecture() (config.Arcitecture, error) {
	arcitectures := []config.Arcitecture{config.Hexagonal, config.Module, config.MVC}
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

func createConfigFile(config *config.Config) error {
	rootPath, _ := filepath.Abs(".")
	configPath := filepath.Join(rootPath, _CONFIG_FILE)
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	err = utils.WriteFile(configPath, string(yamlBytes))
	if err != nil {
		return err
	}
	return nil
}
