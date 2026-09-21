package main

import (
	"os"

	"github.com/sboy99/go-vault/internal/bootstrap"
	"github.com/sboy99/go-vault/internal/config"
	transportcli "github.com/sboy99/go-vault/internal/transport/cli"
	"github.com/sboy99/go-vault/pkg/logger"
)

func main() {
	logger.Init(logger.INFO)

	if needsConfig() {
		if err := config.Load(); err != nil {
			logger.Error("%s", err.Error())
			os.Exit(1)
		}
	} else {
		config.LoadOptional()
	}

	cfg := config.GetConfig()
	c, err := bootstrap.New(cfg)
	if err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := c.Close(); err != nil {
			logger.Error("%s", err.Error())
		}
	}()

	root := transportcli.NewRootCommand(transportcli.Deps{
		Backup: c.Backup,
		Jobs:   c.Jobs,
		Setup:  c.Setup,
	})
	if err := root.Execute(); err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
}

func needsConfig() bool {
	if len(os.Args) < 2 {
		return false
	}
	switch os.Args[1] {
	case "setup", "help", "completion", "--help", "-h", "--version", "-v":
		return false
	default:
		return true
	}
}
