package cli

import (
	"fmt"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/sboy99/go-vault/internal/domain"
	"github.com/sboy99/go-vault/internal/ui"
	"github.com/spf13/cobra"
)

func newSetupCommand(setup *app.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Interactive setup of database and storage config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings, err := collectSetupAnswers()
			if err != nil {
				return err
			}
			if err := setup.SaveSettings(settings); err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			return nil
		},
	}
}

func collectSetupAnswers() (domain.Settings, error) {
	var settings domain.Settings

	dbType, err := ui.DisplaySelectDatabaseTypePrompt()
	if err != nil {
		return settings, fmt.Errorf("select database type: %w", err)
	}
	settings.DatabaseType = domain.DatabaseType(dbType)

	dbHost, err := ui.DisplayInputDatabaseHostPrompt()
	if err != nil {
		return settings, fmt.Errorf("input database host: %w", err)
	}
	settings.DBHost = dbHost

	dbPort, err := ui.DisplayInputDatabasePortPrompt(dbType)
	if err != nil {
		return settings, fmt.Errorf("input database port: %w", err)
	}
	settings.DBPort = dbPort

	dbName, err := ui.DisplayInputDatabaseNamePrompt(dbType)
	if err != nil {
		return settings, fmt.Errorf("input database name: %w", err)
	}
	settings.DBName = dbName

	dbUser, err := ui.DisplayInputDatabaseUsernamePrompt(dbType)
	if err != nil {
		return settings, fmt.Errorf("input database username: %w", err)
	}
	settings.DBUsername = dbUser

	dbPass, err := ui.DisplayInputDatabasePasswordPrompt()
	if err != nil {
		return settings, fmt.Errorf("input database password: %w", err)
	}
	settings.DBPassword = dbPass
	settings.DBSSLMode = "require"

	storageType, err := ui.DisplaySelectStorageTypePrompt()
	if err != nil {
		return settings, fmt.Errorf("select storage type: %w", err)
	}
	settings.StorageType = domain.StorageType(storageType)

	if storageType != "CLOUD" {
		return settings, nil
	}

	cloudType, err := ui.DisplaySelectCloudTypePrompt()
	if err != nil {
		return settings, fmt.Errorf("select cloud type: %w", err)
	}
	settings.CloudType = string(cloudType)

	if cloudType != "AWS" {
		return settings, nil
	}

	awsRegion, err := ui.DisplaySelectAWSRegionPrompt()
	if err != nil {
		return settings, fmt.Errorf("select aws region: %w", err)
	}
	settings.AWSRegion = awsRegion

	awsBucketName, err := ui.DisplayInputAWSBucketNamePrompt()
	if err != nil {
		return settings, fmt.Errorf("input aws bucket name: %w", err)
	}
	settings.AWSBucketName = awsBucketName

	awsAccessKeyId, err := ui.DisplayInputAWSAccessKeyIdPrompt()
	if err != nil {
		return settings, fmt.Errorf("input aws access key id: %w", err)
	}
	settings.AWSAccessKeyID = awsAccessKeyId

	awsAccessKeySecret, err := ui.DisplayInputAWSAccessKeySecretPrompt()
	if err != nil {
		return settings, fmt.Errorf("input aws access key secret: %w", err)
	}
	settings.AWSAccessSecret = awsAccessKeySecret

	awsEndpoint, err := ui.DisplayInputAWSEndpointPrompt()
	if err != nil {
		return settings, fmt.Errorf("input aws endpoint: %w", err)
	}
	settings.AWSEndpoint = awsEndpoint
	return settings, nil
}
