package setup

import (
	"fmt"

	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/ui"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (c *ConfigService) SetupConfig() error {
	cfg := config.GetConfig()
	if err := c.promptDatabase(cfg); err != nil {
		return err
	}
	if err := c.promptStorage(cfg); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (c *ConfigService) promptDatabase(cfg *config.Config) error {
	dbType, err := ui.DisplaySelectDatabaseTypePrompt()
	if err != nil {
		return fmt.Errorf("select database type: %w", err)
	}
	cfg.DB.Type = dbType

	dbHost, err := ui.DisplayInputDatabaseHostPrompt()
	if err != nil {
		return fmt.Errorf("input database host: %w", err)
	}
	cfg.DB.Host = dbHost

	dbPort, err := ui.DisplayInputDatabasePortPrompt(dbType)
	if err != nil {
		return fmt.Errorf("input database port: %w", err)
	}
	cfg.DB.Port = dbPort

	dbName, err := ui.DisplayInputDatabaseNamePrompt(dbType)
	if err != nil {
		return fmt.Errorf("input database name: %w", err)
	}
	cfg.DB.Name = dbName

	dbUser, err := ui.DisplayInputDatabaseUsernamePrompt(dbType)
	if err != nil {
		return fmt.Errorf("input database username: %w", err)
	}
	cfg.DB.Username = dbUser

	dbPass, err := ui.DisplayInputDatabasePasswordPrompt()
	if err != nil {
		return fmt.Errorf("input database password: %w", err)
	}
	cfg.DB.Password = dbPass
	cfg.DB.SSLMode = "require"
	return nil
}

func (c *ConfigService) promptStorage(cfg *config.Config) error {
	storageType, err := ui.DisplaySelectStorageTypePrompt()
	if err != nil {
		return fmt.Errorf("select storage type: %w", err)
	}
	cfg.Storage.Type = storageType

	if storageType != config.CLOUD {
		return nil
	}

	cloudType, err := ui.DisplaySelectCloudTypePrompt()
	if err != nil {
		return fmt.Errorf("select cloud type: %w", err)
	}
	cfg.Storage.Cloud.Type = cloudType

	if cloudType == config.AWS {
		return c.promptAWS(cfg)
	}
	return nil
}

func (c *ConfigService) promptAWS(cfg *config.Config) error {
	awsRegion, err := ui.DisplaySelectAWSRegionPrompt()
	if err != nil {
		return fmt.Errorf("select aws region: %w", err)
	}
	cfg.Storage.Cloud.AWS.Region = awsRegion

	awsBucketName, err := ui.DisplayInputAWSBucketNamePrompt()
	if err != nil {
		return fmt.Errorf("input aws bucket name: %w", err)
	}
	cfg.Storage.Cloud.AWS.BucketName = awsBucketName

	awsAccessKeyId, err := ui.DisplayInputAWSAccessKeyIdPrompt()
	if err != nil {
		return fmt.Errorf("input aws access key id: %w", err)
	}
	cfg.Storage.Cloud.AWS.AccessKeyId = awsAccessKeyId

	awsAccessKeySecret, err := ui.DisplayInputAWSAccessKeySecretPrompt()
	if err != nil {
		return fmt.Errorf("input aws access key secret: %w", err)
	}
	cfg.Storage.Cloud.AWS.AccessKeySecret = awsAccessKeySecret

	awsEndpoint, err := ui.DisplayInputAWSEndpointPrompt()
	if err != nil {
		return fmt.Errorf("input aws endpoint: %w", err)
	}
	cfg.Storage.Cloud.AWS.Endpoint = awsEndpoint
	return nil
}
