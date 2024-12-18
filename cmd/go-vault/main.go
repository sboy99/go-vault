package main

import (
	"github.com/sboy99/go-vault/internal/services"
	"github.com/sboy99/go-vault/internal/strategies"
	"github.com/sboy99/go-vault/pkg/cmd"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	// Logger
	logger.Init(logger.DEBUG)
	// CMD
	if err := cmd.Execute(); err != nil {
		logger.Panic("Error executing CMD... %v", err)
	}
	// Service
	service := services.NewDatabaseService()
	if err := service.Backup(strategies.POSTGRES, "pluto", "localhost", 5432, "postgres", "postgres"); err != nil {
		logger.Panic("Error connecting to database... %v", err)
	}

}
