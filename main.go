package main

import (
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/cmd"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	// Logger
	logger.Init(logger.DEBUG)
	// User Config //
	if err := config.Load(); err != nil {
		logger.Error("%s", err.Error())
		return
	}

	cfg := config.GetConfig()
	logger.Debug("Config %v", cfg.App.Name)
	// Logger //
	logger.Init(logger.DEBUG)
	// CMD //
	if err := cmd.Execute(); err != nil {
		logger.Panic("Error executing CMD\n%v", err)
	}
}
