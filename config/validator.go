package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("app.name", "go-vault")
	viper.SetDefault("app.version", "0.0.1")

	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.dest", "./generated")
}

func validateConfig() error {
	var errorList []string

	if viper.GetString("app.name") == "" {
		errorList = append(errorList, "Missing app name")
	}
	if viper.GetString("app.version") == "" {
		errorList = append(errorList, "Missing app version")
	}

	if viper.GetString("db.name") == "" {
		errorList = append(errorList, "Missing db name")
	}
	if viper.GetString("db.host") == "" {
		errorList = append(errorList, "Missing db host")
	}
	if viper.GetInt("db.port") == 0 {
		errorList = append(errorList, "Missing db port")
	}
	if viper.GetString("db.username") == "" {
		errorList = append(errorList, "Missing db username")
	}
	if viper.GetString("db.password") == "" {
		errorList = append(errorList, "Missing db password")
	}

	if viper.GetString("storage.type") == "" {
		errorList = append(errorList, "Missing storage type")
	}
	if viper.GetString("storage.dest") == "" {
		errorList = append(errorList, "Missing storage dest")
	}

	if len(errorList) > 0 {
		return fmt.Errorf("Config Error\n%v", strings.Join(errorList, "\n"))
	}

	return nil
}
