package config

import (
	"fmt"
	"time"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/spf13/viper"
)

// ToConnParams maps DB config to domain connection params.
func (c *Config) ToConnParams() domain.ConnParams {
	return domain.ConnParams{
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		Name:     c.DB.Name,
		Username: c.DB.Username,
		Password: c.DB.Password,
		SSLMode:  c.DB.SSLMode,
	}
}

// ToBackupConfig maps runtime config into the narrow app.BackupConfig surface.
func (c *Config) ToBackupConfig() (app.BackupConfig, error) {
	loc, err := time.LoadLocation(c.Schedule.Timezone)
	if err != nil {
		return app.BackupConfig{}, fmt.Errorf("load timezone %q: %w", c.Schedule.Timezone, err)
	}
	return app.BackupConfig{
		DBName:         c.DB.Name,
		DatabaseType:   domain.DatabaseType(c.DB.Type),
		StorageType:    domain.StorageType(c.Storage.Type),
		Conn:           c.ToConnParams(),
		TempDir:        c.Runtime.TempDir,
		DumpTimeout:    c.Runtime.DumpTimeout,
		RestoreTimeout: c.Runtime.RestoreTimeout,
		RestoreJobs:    c.Runtime.RestoreJobs,
		Retention: domain.RetentionPolicy{
			Daily:   c.Retention.Daily,
			Weekly:  c.Retention.Weekly,
			Monthly: c.Retention.Monthly,
		},
		Location: loc,
	}, nil
}

// SettingsStore persists domain.Settings via viper.
type SettingsStore struct{}

var _ domain.SettingsStore = (*SettingsStore)(nil)

func NewSettingsStore() *SettingsStore {
	return &SettingsStore{}
}

func (SettingsStore) Save(s domain.Settings) error {
	cfg := GetConfig()
	cfg.DB.Type = DatabaseEnum(s.DatabaseType)
	cfg.DB.Host = s.DBHost
	cfg.DB.Port = s.DBPort
	cfg.DB.Name = s.DBName
	cfg.DB.Username = s.DBUsername
	cfg.DB.Password = s.DBPassword
	cfg.DB.SSLMode = s.DBSSLMode
	cfg.Storage.Type = StorageEnum(s.StorageType)
	if s.StorageDest != "" {
		cfg.Storage.Dest = s.StorageDest
	}
	cfg.Storage.Cloud.Type = CloudEnum(s.CloudType)
	cfg.Storage.Cloud.AWS.Region = s.AWSRegion
	cfg.Storage.Cloud.AWS.BucketName = s.AWSBucketName
	cfg.Storage.Cloud.AWS.AccessKeyId = s.AWSAccessKeyID
	cfg.Storage.Cloud.AWS.AccessKeySecret = s.AWSAccessSecret
	cfg.Storage.Cloud.AWS.Endpoint = s.AWSEndpoint
	return Save(cfg)
}

// ApplySettingsToViper is used by tests that need to inspect viper keys.
func ApplySettingsToViper(s domain.Settings) {
	viper.Set("db.type", string(s.DatabaseType))
	viper.Set("db.host", s.DBHost)
	viper.Set("db.port", s.DBPort)
	viper.Set("db.name", s.DBName)
	viper.Set("db.username", s.DBUsername)
	viper.Set("db.password", s.DBPassword)
	viper.Set("db.sslmode", s.DBSSLMode)
	viper.Set("storage.type", string(s.StorageType))
}
