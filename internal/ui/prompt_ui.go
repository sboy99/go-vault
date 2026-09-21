package ui

import (
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-vault/internal/config"
)

var awsRegions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"ca-central-1",
	"sa-east-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"eu-central-1",
	"eu-central-2",
	"eu-north-1",
	"eu-south-1",
	"eu-south-2",
	"me-south-1",
	"me-central-1",
	"il-central-1",
	"af-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-southeast-3",
	"ap-southeast-4",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-south-2",
	"ap-east-1",
	"ap-southeast-5",
}

func inputPrompt(label, def string, mask rune) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: def,
		Mask:    mask,
	}
	return prompt.Run()
}

func selectPrompt[T ~string](label string, items []T) (T, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		var zero T
		return zero, err
	}
	return T(result), nil
}

func DisplaySelectDatabaseTypePrompt() (config.DatabaseEnum, error) {
	return selectPrompt("Select DB", []config.DatabaseEnum{config.POSTGRESQL})
}

func DisplayInputDatabaseNamePrompt(dbType config.DatabaseEnum) (string, error) {
	return inputPrompt("Enter DB Name", getDatabaseName(dbType), 0)
}

func DisplayInputDatabaseHostPrompt() (string, error) {
	return inputPrompt("Enter DB Host", "localhost", 0)
}

func DisplayInputDatabasePortPrompt(dbType config.DatabaseEnum) (int, error) {
	prompt := promptui.Prompt{
		Label:    "Enter DB Port",
		Validate: portValidator,
		Default:  getDatabasePort(dbType),
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result)
}

func DisplayInputDatabaseUsernamePrompt(dbType config.DatabaseEnum) (string, error) {
	return inputPrompt("Enter DB Username", getDatabaseUser(dbType), 0)
}

func DisplayInputDatabasePasswordPrompt() (string, error) {
	return inputPrompt("Enter DB Password", "", '*')
}

func DisplaySelectStorageTypePrompt() (config.StorageEnum, error) {
	return selectPrompt("Select Storage", []config.StorageEnum{config.LOCAL, config.CLOUD})
}

func DisplaySelectCloudTypePrompt() (config.CloudEnum, error) {
	return selectPrompt("Select Cloud", []config.CloudEnum{config.AWS})
}

func DisplaySelectAWSRegionPrompt() (string, error) {
	return selectPrompt("Enter AWS Region", awsRegions)
}

func DisplayInputAWSBucketNamePrompt() (string, error) {
	return inputPrompt("Enter AWS Bucket Name", "", 0)
}

func DisplayInputAWSAccessKeyIdPrompt() (string, error) {
	return inputPrompt("Enter AWS Access Key ID", "", 0)
}

func DisplayInputAWSAccessKeySecretPrompt() (string, error) {
	return inputPrompt("Enter AWS Access Key Secret", "", '*')
}

func DisplayInputAWSEndpointPrompt() (string, error) {
	return inputPrompt("Enter AWS Endpoint", "default", 0)
}

func DisplayRestoreConfirmPrompt() (string, error) {
	return inputPrompt("Type the database name to confirm restore", "", 0)
}
