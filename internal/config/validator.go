package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/sboy99/go-vault/internal/version"
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("app.name", "go-vault")
	viper.SetDefault("app.version", version.Version)

	viper.SetDefault("db.type", string(POSTGRESQL))
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.sslmode", "require")

	viper.SetDefault("storage.type", string(LOCAL))
	viper.SetDefault("storage.dest", "./backups")
	viper.SetDefault("storage.cloud.aws.endpoint", "default")

	viper.SetDefault("schedule.cron", "0 2 * * *")
	viper.SetDefault("schedule.timezone", "UTC")

	viper.SetDefault("retention.daily", 7)
	viper.SetDefault("retention.weekly", 4)
	viper.SetDefault("retention.monthly", 12)

	viper.SetDefault("api.addr", ":8080")
	viper.SetDefault("api.token", "")

	viper.SetDefault("runtime.temp_dir", "/tmp/go-vault")
	viper.SetDefault("runtime.dump_timeout", 2*time.Hour)
	viper.SetDefault("runtime.restore_timeout", 2*time.Hour)
	viper.SetDefault("runtime.shutdown_timeout", 5*time.Minute)
	viper.SetDefault("runtime.meta_db_path", "./go-vault.db")
	viper.SetDefault("runtime.alert_webhook", "")
	viper.SetDefault("runtime.restore_jobs", 4)
}

func validateConfig() error {
	var errorList []string

	if viper.GetString("app.name") == "" {
		errorList = append(errorList, "missing app.name")
	}
	if viper.GetString("db.name") == "" {
		errorList = append(errorList, "missing db.name")
	}
	if viper.GetString("db.host") == "" {
		errorList = append(errorList, "missing db.host")
	}
	if viper.GetInt("db.port") == 0 {
		errorList = append(errorList, "missing db.port")
	}
	if viper.GetString("db.username") == "" {
		errorList = append(errorList, "missing db.username")
	}
	if viper.GetString("db.password") == "" {
		errorList = append(errorList, "missing db.password")
	}
	if viper.GetString("db.sslmode") == "" {
		errorList = append(errorList, "missing db.sslmode")
	}

	dbType := strings.ToUpper(viper.GetString("db.type"))
	if dbType != "" && dbType != string(POSTGRESQL) {
		errorList = append(errorList, fmt.Sprintf("unsupported db.type %q (only POSTGRESQL is implemented)", dbType))
	}

	storageType := strings.ToUpper(viper.GetString("storage.type"))
	if storageType == "" {
		errorList = append(errorList, "missing storage.type")
	}
	switch StorageEnum(storageType) {
	case LOCAL:
		if viper.GetString("storage.dest") == "" {
			errorList = append(errorList, "missing storage.dest for LOCAL storage")
		}
	case CLOUD:
		cloudType := strings.ToUpper(viper.GetString("storage.cloud.type"))
		if cloudType != string(AWS) {
			errorList = append(errorList, "unsupported storage.cloud.type (only AWS is implemented)")
		}
		if viper.GetString("storage.cloud.aws.region") == "" {
			errorList = append(errorList, "missing storage.cloud.aws.region")
		}
		if viper.GetString("storage.cloud.aws.bucket_name") == "" {
			errorList = append(errorList, "missing storage.cloud.aws.bucket_name")
		}
		if viper.GetString("storage.cloud.aws.access_key_id") == "" {
			errorList = append(errorList, "missing storage.cloud.aws.access_key_id")
		}
		if viper.GetString("storage.cloud.aws.access_key_secret") == "" {
			errorList = append(errorList, "missing storage.cloud.aws.access_key_secret")
		}
	default:
		if storageType != "" {
			errorList = append(errorList, fmt.Sprintf("unsupported storage.type %q", storageType))
		}
	}

	if viper.GetString("schedule.cron") == "" {
		errorList = append(errorList, "missing schedule.cron")
	}
	if viper.GetString("schedule.timezone") == "" {
		errorList = append(errorList, "missing schedule.timezone")
	}
	if viper.GetInt("retention.daily") < 1 {
		errorList = append(errorList, "retention.daily must be >= 1")
	}
	if viper.GetInt("retention.weekly") < 1 {
		errorList = append(errorList, "retention.weekly must be >= 1")
	}
	if viper.GetInt("retention.monthly") < 1 {
		errorList = append(errorList, "retention.monthly must be >= 1")
	}
	if viper.GetString("api.addr") == "" {
		errorList = append(errorList, "missing api.addr")
	}
	if viper.GetDuration("runtime.dump_timeout") <= 0 {
		errorList = append(errorList, "runtime.dump_timeout must be > 0")
	}
	if viper.GetDuration("runtime.restore_timeout") <= 0 {
		errorList = append(errorList, "runtime.restore_timeout must be > 0")
	}
	if viper.GetInt("runtime.restore_jobs") < 1 {
		errorList = append(errorList, "runtime.restore_jobs must be >= 1")
	}

	if len(errorList) > 0 {
		return fmt.Errorf("config error:\n%s", strings.Join(errorList, "\n"))
	}
	return nil
}
