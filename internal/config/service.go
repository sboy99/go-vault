package config

import (
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/sboy99/go-vault/pkg/logger"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// TODO: Return error
func (c *ConfigService) SetupConfig() {
	// Config //
	cfg := config.GetConfig()

	// Promt input //
	dbType, err := ui.DisplaySelectDatabaseTypePrompt()
	if err != nil {
		logger.Error("Failed to display select db prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Type = dbType

	dbHost, err := ui.DisplayInputDatabaseHostPrompt()
	if err != nil {
		logger.Error("Failed to display input db host prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Host = dbHost

	dbPort, err := ui.DisplayInputDatabasePortPrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db port prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Port = dbPort

	dbName, err := ui.DisplayInputDatabaseNamePrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db name prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Name = dbName

	dbUser, err := ui.DisplayInputDatabaseUsernamePrompt(dbType)
	if err != nil {
		logger.Error("Failed to display input db username prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Username = dbUser

	dbPass, err := ui.DisplayInputDatabasePasswordPrompt()
	if err != nil {
		logger.Error("Failed to display input db password prompt\nDetails: %v", err)
		return
	}
	cfg.DB.Password = dbPass

	storageType, err := ui.DisplaySelectStorageTypePrompt()
	if err != nil {
		logger.Error("Failed to display select storage prompt\nDetails: %v", err)
		return
	}
	cfg.Storage.Type = storageType

	if storageType == config.CLOUD {
		cloudType, err := ui.DisplaySelectCloudTypePrompt()
		if err != nil {
			logger.Error("Failed to display select cloud prompt\nDetails: %v", err)
			return
		}
		cfg.Storage.Cloud.Type = cloudType

		if cloudType == config.AWS {
			awsRegion, err := ui.DisplaySelectAWSRegionPrompt()
			if err != nil {
				logger.Error("Failed to display input aws region prompt\nDetails: %v", err)
				return
			}
			cfg.Storage.Cloud.AWS.Region = awsRegion

			awsBucketName, err := ui.DisplayInputAWSBucketNamePrompt()
			if err != nil {
				logger.Error("Failed to display input aws bucket name prompt\nDetails: %v", err)
				return
			}
			cfg.Storage.Cloud.AWS.BucketName = awsBucketName

			awsAccessKeyId, err := ui.DisplayInputAWSAccessKeyIdPrompt()
			if err != nil {
				logger.Error("Failed to display input aws access key id prompt\nDetails: %v", err)
				return
			}
			cfg.Storage.Cloud.AWS.AccessKeyId = awsAccessKeyId

			awsAccessKeySecret, err := ui.DisplayInputAWSAccessKeySecretPrompt()
			if err != nil {
				logger.Error("Failed to display input aws access key secret prompt\nDetails: %v", err)
				return
			}
			cfg.Storage.Cloud.AWS.AccessKeySecret = awsAccessKeySecret

			awsEndpoint, err := ui.DisplayInputAWSEndpointPrompt()
			if err != nil {
				logger.Error("Failed to display input aws endpoint prompt\nDetails: %v", err)
				return
			}
			cfg.Storage.Cloud.AWS.Endpoint = awsEndpoint
		}
	}

	if err = config.Save(cfg); err != nil {
		logger.Error("Failed to save config\nDetails: %v", err)
	}
}
