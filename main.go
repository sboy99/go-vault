package main

import (
	"os"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/cmd"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	// Logger
	logger.Init(logger.DEBUG)

	logger.Info("args %v", os.Args)

	// User Config //
	if len(os.Args) > 1 && os.Args[1] != "setup" {
		if err := config.Load(); err != nil {
			logger.Error("%s", err.Error())
			return
		}
	}

	// Cmd //
	if err := cmd.Execute(); err != nil {
		logger.Error("%s", err.Error())
		return
	}

}
