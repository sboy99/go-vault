package config

import (
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/sboy99/go-vault/pkg/logger"
)

func SetupConfig() {
	// Promt input //
	dbType, err := ui.DisplaySelectDatabaseTypePrompt()
	if err != nil {
		logger.Error("Failed to display select db prompt\nDetails: %v", err)
	}
	dbHost, err := ui.DisplayInputDatabaseHostPrompt()
	if err != nil {
		logger.Error("Failed to display input db host prompt\nDetails: %v", err)
	}
	dbPort, err := ui.DisplayInputDatabasePortPrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db port prompt\nDetails: %v", err)
	}
	dbName, err := ui.DisplayInputDatabaseNamePrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db name prompt\nDetails: %v", err)
	}
	dbUser, err := ui.DisplayInputDatabaseUsernamePrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db username prompt\nDetails: %v", err)
	}
	dbPass, err := ui.DisplayInputDatabasePasswordPrompt()
	if err != nil {
		logger.Error("Failed to display input db password prompt\nDetails: %v", err)
	}

	// Save config //
	cfg := config.GetConfig()
	cfg.DB.Type = dbType
	cfg.DB.Host = dbHost
	cfg.DB.Port = dbPort
	cfg.DB.Name = dbName
	cfg.DB.Username = dbUser
	cfg.DB.Password = dbPass
	if err = config.Save(cfg); err != nil {
		logger.Error("Failed to save config\nDetails: %v", err)
	}
}
