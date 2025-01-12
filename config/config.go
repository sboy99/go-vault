package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App     app
	DB      database
	Storage storage
}

type app struct {
	Name    string
	Version string
}

type DatabaseEnum string
type database struct {
	Name     string
	Type     DatabaseEnum
	Host     string
	Port     int
	Username string
	Password string
}

type StorageEnum string
type storage struct {
	Type  StorageEnum
	Dest  string //TODO: Move Dest to LocalStorage struct
	Cloud cloudStorage
}

type CloudEnum string
type cloudStorage struct {
	Type CloudEnum
	AWS  _AWSCloudStorage
}

type _AWSCloudStorage struct {
	Region          string
	BucketName      string
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

const (
	POSTGRESQL DatabaseEnum = "POSTGRESQL"
	MYSQL      DatabaseEnum = "MYSQL"
	MONGODB    DatabaseEnum = "MONGODB"
)

const (
	LOCAL StorageEnum = "local"
	CLOUD StorageEnum = "cloud"
)

const (
	GCP CloudEnum = "gcp"
	AWS CloudEnum = "aws"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	setDefaults()
}

func Load() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("No config is setup please run `go-vault setup`")
	}
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
	// Storage //
	viper.Set("storage.type", config.Storage.Type)
	viper.Set("storage.dest", config.Storage.Dest)
	// Cloud Storage //
	viper.Set("storage.cloud.type", config.Storage.Cloud.Type)
	// AWS Cloud Storage //
	viper.Set("storage.cloud.aws.region", config.Storage.Cloud.AWS.Region)
	viper.Set("storage.cloud.aws.bucket_name", config.Storage.Cloud.AWS.BucketName)
	viper.Set("storage.cloud.aws.access_key_id", config.Storage.Cloud.AWS.AccessKeyId)
	viper.Set("storage.cloud.aws.access_key_secret", config.Storage.Cloud.AWS.AccessKeySecret)
	viper.Set("storage.cloud.aws.endpoint", config.Storage.Cloud.AWS.Endpoint)

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
			Type:     DatabaseEnum(viper.GetString("db.type")),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		},
		Storage: storage{
			Type: StorageEnum(viper.GetString("storage.type")),
			Dest: viper.GetString("storage.dest"),
			Cloud: cloudStorage{
				Type: CloudEnum(viper.GetString("storage.cloud.type")),
				AWS: _AWSCloudStorage{
					Region:          viper.GetString("storage.cloud.aws.region"),
					BucketName:      viper.GetString("storage.cloud.aws.bucket_name"),
					AccessKeyId:     viper.GetString("storage.cloud.aws.access_key_id"),
					AccessKeySecret: viper.GetString("storage.cloud.aws.access_key_secret"),
					Endpoint:        viper.GetString("storage.cloud.aws.endpoint"),
				},
			},
		},
	}
}
