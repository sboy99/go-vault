package ui

import (
	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/utils"
)

func DisplaySelectDatabaseTypePrompt() (config.DatabaseEnum, error) {
	// Prompt for DB selection //
	prompt := promptui.Select{
		Label: "Select DB",
		Items: []config.DatabaseEnum{config.POSTGRESQL, config.MYSQL, config.MONGODB},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.DatabaseEnum(result), nil
}

func DisplayInputDatabaseNamePrompt(dbType config.DatabaseEnum) (string, error) {
	// Prompt for DB name //
	prompt := promptui.Prompt{
		Label:   "Enter DB Name",
		Default: getDatabaseName(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabaseHostPrompt() (string, error) {
	// Prompt for DB host //
	prompt := promptui.Prompt{
		Label:   "Enter DB Host",
		Default: "localhost",
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabasePortPrompt(dbType config.DatabaseEnum) (int, error) {
	// Prompt for DB port //
	prompt := promptui.Prompt{
		Label:    "Enter DB Port",
		Validate: portValidator,
		Default:  getDatabasePort(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	resultInt, err := utils.ParseInt(result)
	if err != nil {
		return 0, err
	}
	return resultInt, nil
}

func DisplayInputDatabaseUsernamePrompt(dbType config.DatabaseEnum) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Enter DB Username",
		Default: getDatabaseUser(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabasePasswordPrompt() (string, error) {
	// Prompt for DB password //
	prompt := promptui.Prompt{
		Label: "Enter DB Password",
		Mask:  '*',
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplaySelectStorageTypePrompt() (config.StorageEnum, error) {
	// Prompt for Storage selection //
	prompt := promptui.Select{
		Label: "Select Storage",
		Items: []config.StorageEnum{config.LOCAL, config.CLOUD},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.StorageEnum(result), nil
}
