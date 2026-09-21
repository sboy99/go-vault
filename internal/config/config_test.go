package config_test

import (
	"testing"

	"github.com/sboy99/go-vault/internal/config"
	"github.com/spf13/viper"
)

func TestGetConfigUppercasesEnums(t *testing.T) {
	viper.Set("db.type", "postgresql")
	viper.Set("storage.type", "local")
	cfg := config.GetConfig()
	if cfg.DB.Type != config.POSTGRESQL {
		t.Fatalf("db type %q", cfg.DB.Type)
	}
	if cfg.Storage.Type != config.LOCAL {
		t.Fatalf("storage type %q", cfg.Storage.Type)
	}
}

func TestRetentionDefaults(t *testing.T) {
	config.LoadOptional()
	cfg := config.GetConfig()
	if cfg.Retention.Daily < 1 || cfg.Retention.Weekly < 1 || cfg.Retention.Monthly < 1 {
		t.Fatalf("retention defaults invalid: %+v", cfg.Retention)
	}
}
