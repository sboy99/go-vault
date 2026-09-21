package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       App
	DB        Database
	Storage   Storage
	Schedule  Schedule
	Retention Retention
	API       API
	Runtime   Runtime
}

type App struct {
	Name    string
	Version string
}

type DatabaseEnum string

const (
	POSTGRESQL DatabaseEnum = "POSTGRESQL"
)

type Database struct {
	Name     string
	Type     DatabaseEnum
	Host     string
	Port     int
	Username string
	Password string
	SSLMode  string
}

type StorageEnum string

const (
	LOCAL StorageEnum = "LOCAL"
	CLOUD StorageEnum = "CLOUD"
)

type Storage struct {
	Type  StorageEnum
	Dest  string
	Cloud CloudStorage
}

type CloudEnum string

const (
	AWS CloudEnum = "AWS"
)

type CloudStorage struct {
	Type CloudEnum
	AWS  AWSCloudStorage
}

type AWSCloudStorage struct {
	Region          string
	BucketName      string
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

type Schedule struct {
	Cron     string
	Timezone string
}

type Retention struct {
	Daily   int
	Weekly  int
	Monthly int
}

type API struct {
	Addr  string
	Token string
}

type Runtime struct {
	TempDir         string
	DumpTimeout     time.Duration
	RestoreTimeout  time.Duration
	ShutdownTimeout time.Duration
	MetaDBPath      string
	AlertWebhook    string
	RestoreJobs     int
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetEnvPrefix("GO_VAULT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	setDefaults()
}

// Load reads optional YAML, overlays env vars, resolves _FILE secrets, and validates.
func Load() error {
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return fmt.Errorf("read config: %w", err)
		}
	}
	bindFileSecrets()
	if err := validateConfig(); err != nil {
		return err
	}
	return nil
}

// LoadOptional loads config without requiring DB credentials (used by setup).
func LoadOptional() {
	_ = viper.ReadInConfig()
	bindFileSecrets()
}

func bindFileSecrets() {
	resolveFileEnv("db.password", "GO_VAULT_DB_PASSWORD_FILE")
	resolveFileEnv("storage.cloud.aws.access_key_secret", "GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_SECRET_FILE")
	resolveFileEnv("api.token", "GO_VAULT_API_TOKEN_FILE")
}

func resolveFileEnv(configKey, fileEnv string) {
	path := os.Getenv(fileEnv)
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	viper.Set(configKey, strings.TrimSpace(string(data)))
}

func Save(cfg *Config) error {
	viper.Set("app.name", cfg.App.Name)
	viper.Set("app.version", cfg.App.Version)
	viper.Set("db.name", cfg.DB.Name)
	viper.Set("db.type", cfg.DB.Type)
	viper.Set("db.host", cfg.DB.Host)
	viper.Set("db.port", cfg.DB.Port)
	viper.Set("db.username", cfg.DB.Username)
	viper.Set("db.password", cfg.DB.Password)
	viper.Set("db.sslmode", cfg.DB.SSLMode)
	viper.Set("storage.type", cfg.Storage.Type)
	viper.Set("storage.dest", cfg.Storage.Dest)
	viper.Set("storage.cloud.type", cfg.Storage.Cloud.Type)
	viper.Set("storage.cloud.aws.region", cfg.Storage.Cloud.AWS.Region)
	viper.Set("storage.cloud.aws.bucket_name", cfg.Storage.Cloud.AWS.BucketName)
	viper.Set("storage.cloud.aws.access_key_id", cfg.Storage.Cloud.AWS.AccessKeyId)
	viper.Set("storage.cloud.aws.access_key_secret", cfg.Storage.Cloud.AWS.AccessKeySecret)
	viper.Set("storage.cloud.aws.endpoint", cfg.Storage.Cloud.AWS.Endpoint)
	viper.Set("schedule.cron", cfg.Schedule.Cron)
	viper.Set("schedule.timezone", cfg.Schedule.Timezone)
	viper.Set("retention.daily", cfg.Retention.Daily)
	viper.Set("retention.weekly", cfg.Retention.Weekly)
	viper.Set("retention.monthly", cfg.Retention.Monthly)
	viper.Set("api.addr", cfg.API.Addr)
	viper.Set("api.token", cfg.API.Token)

	if err := viper.SafeWriteConfig(); err != nil {
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

func GetConfig() *Config {
	return &Config{
		App: App{
			Name:    viper.GetString("app.name"),
			Version: viper.GetString("app.version"),
		},
		DB: Database{
			Name:     viper.GetString("db.name"),
			Type:     DatabaseEnum(strings.ToUpper(viper.GetString("db.type"))),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
		Storage: Storage{
			Type: StorageEnum(strings.ToUpper(viper.GetString("storage.type"))),
			Dest: viper.GetString("storage.dest"),
			Cloud: CloudStorage{
				Type: CloudEnum(strings.ToUpper(viper.GetString("storage.cloud.type"))),
				AWS: AWSCloudStorage{
					Region:          viper.GetString("storage.cloud.aws.region"),
					BucketName:      viper.GetString("storage.cloud.aws.bucket_name"),
					AccessKeyId:     viper.GetString("storage.cloud.aws.access_key_id"),
					AccessKeySecret: viper.GetString("storage.cloud.aws.access_key_secret"),
					Endpoint:        viper.GetString("storage.cloud.aws.endpoint"),
				},
			},
		},
		Schedule: Schedule{
			Cron:     viper.GetString("schedule.cron"),
			Timezone: viper.GetString("schedule.timezone"),
		},
		Retention: Retention{
			Daily:   viper.GetInt("retention.daily"),
			Weekly:  viper.GetInt("retention.weekly"),
			Monthly: viper.GetInt("retention.monthly"),
		},
		API: API{
			Addr:  viper.GetString("api.addr"),
			Token: viper.GetString("api.token"),
		},
		Runtime: Runtime{
			TempDir:         viper.GetString("runtime.temp_dir"),
			DumpTimeout:     viper.GetDuration("runtime.dump_timeout"),
			RestoreTimeout:  viper.GetDuration("runtime.restore_timeout"),
			ShutdownTimeout: viper.GetDuration("runtime.shutdown_timeout"),
			MetaDBPath:      viper.GetString("runtime.meta_db_path"),
			AlertWebhook:    viper.GetString("runtime.alert_webhook"),
			RestoreJobs:     viper.GetInt("runtime.restore_jobs"),
		},
	}
}
