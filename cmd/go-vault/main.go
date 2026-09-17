package main

import (
	"os"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/cmd"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	logger.Init(logger.INFO)

	defer func() {
		if err := meta.Cleanup(); err != nil {
			logger.Error("%s", err.Error())
		}
	}()

	if needsConfig() {
		if err := config.Load(); err != nil {
			logger.Error("%s", err.Error())
			os.Exit(1)
		}
	} else {
		config.LoadOptional()
	}

	cfg := config.GetConfig()
	if err := meta.Init(cfg.Runtime.MetaDBPath); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
}

func needsConfig() bool {
	if len(os.Args) < 2 {
		return false
	}
	switch os.Args[1] {
	case "setup", "help", "completion", "--help", "-h":
		return false
	default:
		return true
	}
}
