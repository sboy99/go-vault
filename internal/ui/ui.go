package ui

import (
	"os/user"

	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/strategies"
)

func DisplaySelectDBPrompt() (strategies.DatabaseEnum, error) {
	// Prompt for DB selection //
	prompt := promptui.Select{
		Label: "Select DB",
		Items: []strategies.DatabaseEnum{strategies.POSTGRES, strategies.MYSQL, strategies.MONGO},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return strategies.DatabaseEnum(result), nil
}

func DisplayDBConfigPrompt() (*config.DBConfig, error) {
	dbConfig := config.NewDBConfig()

	// Prompt for DB name //
	dbNamePromt := promptui.Prompt{
		Label: "Enter DB Name",
	}
	dbName, err := dbNamePromt.Run()
	if err != nil {
		return nil, err
	}
	dbConfig.Name = dbName

	// Prompt for DB host //
	dbHostPromt := promptui.Prompt{
		Label:   "Enter DB Host",
		Default: "localhost",
	}
	dbHost, err := dbHostPromt.Run()
	if err != nil {
		return nil, err
	}
	dbConfig.Host = dbHost

	// Prompt for DB port //
	dbPortPromt := promptui.Prompt{
		Label:    "Enter DB Port",
		Validate: portValidator,
		Default:  "5432",
	}
	dbPort, err := dbPortPromt.Run()
	if err != nil {
		return nil, err
	}
	dbConfig.Port = dbPort

	// Prompt for DB username //
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	dbUsernamePromt := promptui.Prompt{
		Label:   "Enter DB Username",
		Default: currentUser.Username,
	}
	dbUsername, err := dbUsernamePromt.Run()
	if err != nil {
		return nil, err
	}
	dbConfig.User = dbUsername

	// Prompt for DB password //
	dbPasswordPromt := promptui.Prompt{
		Label: "Enter DB Password",
		Mask:  '*',
	}
	dbPassword, err := dbPasswordPromt.Run()
	if err != nil {
		return nil, err
	}
	dbConfig.Password = dbPassword

	return dbConfig, nil
}
