package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App app
	DB  database
}

type app struct {
	Name    string
	Version string
}

type DBEnum string
type database struct {
	Name     string
	Type     DBEnum
	Host     string
	Port     int
	Username string
	Password string
}

const (
	POSTGRESQL DBEnum = "postgresql"
	MYSQL      DBEnum = "mysql"
	MONGODB    DBEnum = "mongodb"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
}

func Load() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("No config is setup please run `go-vault setup`")
	}
	setDefaults()
	if err := validateConfig(); err != nil {
		return err
	}
	return nil
}

func Save(config *Config) error {
	// App //
	viper.Set("app.name", config.App.Name)
	viper.Set("app.version", config.App.Version)
	// DB //
	viper.Set("db.name", config.DB.Name)
	viper.Set("db.type", config.DB.Type)
	viper.Set("db.host", config.DB.Host)
	viper.Set("db.port", config.DB.Port)
	viper.Set("db.username", config.DB.Username)
	viper.Set("db.password", config.DB.Password)

	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &Config{
		App: app{
			Name:    viper.GetString("app.name"),
			Version: viper.GetString("app.version"),
		},
		DB: database{
			Name:     viper.GetString("db.name"),
			Type:     DBEnum(viper.GetString("db.type")),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		},
	}
}
