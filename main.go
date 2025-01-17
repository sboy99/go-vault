package main

import (
	"os"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/cmd"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/pkg/logger"
)

func cleanup() {
	if err := meta.Cleanup(); err != nil {
		logger.Error("%s", err.Error())
	}
}

func main() {
	// Cleanup
	defer cleanup()

	// Logger
	logger.Init(logger.DEBUG)

	// MetaData
	if err := meta.Init(); err != nil {
		logger.Error("%s", err.Error())
		return
	}

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
