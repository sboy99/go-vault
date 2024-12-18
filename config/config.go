package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sboy99/go-vault/pkg/logger"
	"github.com/sboy99/go-vault/pkg/utils"
)

type EnvEnum string

type TConfig struct {
	Env      EnvEnum
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

const (
	PROD EnvEnum = "PROD"
	DEV  EnvEnum = "DEV"
	TEST EnvEnum = "TEST"
)

var (
	Config               = TConfig{}
	requiredEnv          = []string{}
	hasMissingEnv        = false
	invalidNumericEnv    = []string{}
	hasInvalidNumericEnv = false
	invalidBoolEnv       = []string{}
	hasInvalidBoolEnv    = false
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		logger.Panic("No .env file found")
	}

	Config = TConfig{
		Env:      EnvEnum(loadString("Env")),
		Database: loadDatabaseConfig(),
	}

	if hasMissingEnv {
		logger.Panic("Missing environment variables: %v", strings.Join(requiredEnv, ", "))
	}

	if hasInvalidNumericEnv {
		logger.Panic("Invalid numeric environment variables: %v", strings.Join(invalidNumericEnv, ", "))
	}

	if hasInvalidBoolEnv {
		logger.Panic("Invalid boolean environment variables: %v", strings.Join(invalidBoolEnv, ", "))
	}

	logger.Info("Config loaded successfully")
}

func loadDatabaseConfig() DatabaseConfig {
	databaseConfig := DatabaseConfig{
		Name:     loadString("DB_NAME"),
		Host:     loadString("DB_HOST"),
		Port:     loadInt("DB_PORT"),
		User:     loadString("DB_USER"),
		Password: loadString("DB_PASSWORD"),
	}

	return databaseConfig
}

func loadString(key string) string {
	env := os.Getenv(key)
	if env == "" {
		requiredEnv = append(requiredEnv, key)
		hasMissingEnv = true
	}
	return env
}

func loadInt(key string) int {
	env := loadString(key)
	parsedEnv, err := utils.ParseInt(env)
	if err != nil {
		invalidNumericEnv = append(invalidNumericEnv, key)
		hasInvalidNumericEnv = true
	}
	return parsedEnv
}

func loadBool(key string) bool {
	env := loadString(key)
	parsedEnv, err := utils.ParseBool(env)
	if err != nil {
		invalidBoolEnv = append(invalidBoolEnv, key)
		hasInvalidBoolEnv = true
	}
	return parsedEnv
}
