package ui

import (
	"github.com/manifoldco/promptui"
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/utils"
)

func DisplaySelectDatabaseTypePrompt() (config.DatabaseEnum, error) {
	// Prompt for DB selection //
	prompt := promptui.Select{
		Label: "Select DB",
		Items: []config.DatabaseEnum{config.POSTGRESQL, config.MYSQL, config.MONGODB},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.DatabaseEnum(result), nil
}

func DisplayInputDatabaseNamePrompt(dbType config.DatabaseEnum) (string, error) {
	// Prompt for DB name //
	prompt := promptui.Prompt{
		Label:   "Enter DB Name",
		Default: getDatabaseName(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabaseHostPrompt() (string, error) {
	// Prompt for DB host //
	prompt := promptui.Prompt{
		Label:   "Enter DB Host",
		Default: "localhost",
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabasePortPrompt(dbType config.DatabaseEnum) (int, error) {
	// Prompt for DB port //
	prompt := promptui.Prompt{
		Label:    "Enter DB Port",
		Validate: portValidator,
		Default:  getDatabasePort(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	resultInt, err := utils.ParseInt(result)
	if err != nil {
		return 0, err
	}
	return resultInt, nil
}

func DisplayInputDatabaseUsernamePrompt(dbType config.DatabaseEnum) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Enter DB Username",
		Default: getDatabaseUser(dbType),
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputDatabasePasswordPrompt() (string, error) {
	// Prompt for DB password //
	prompt := promptui.Prompt{
		Label: "Enter DB Password",
		Mask:  '*',
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplaySelectStorageTypePrompt() (config.StorageEnum, error) {
	// Prompt for Storage selection //
	prompt := promptui.Select{
		Label: "Select Storage",
		Items: []config.StorageEnum{config.LOCAL, config.CLOUD},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.StorageEnum(result), nil
}

func DisplaySelectCloudTypePrompt() (config.CloudEnum, error) {
	// Prompt for Cloud selection //
	prompt := promptui.Select{
		Label: "Select Cloud",
		Items: []config.CloudEnum{config.AWS, config.GCP},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return config.CloudEnum(result), nil
}

func DisplaySelectAWSRegionPrompt() (string, error) {
	// Prompt for AWS region //
	prompt := promptui.Select{
		Label: "Enter AWS Region",
		Items: []string{
			"us-east-1",      // US East (N. Virginia)
			"us-east-2",      // US East (Ohio)
			"us-west-1",      // US West (N. California)
			"us-west-2",      // US West (Oregon)
			"ca-central-1",   // Canada (Central)
			"sa-east-1",      // South America (São Paulo)
			"eu-west-1",      // Europe (Ireland)
			"eu-west-2",      // Europe (London)
			"eu-west-3",      // Europe (Paris)
			"eu-central-1",   // Europe (Frankfurt)
			"eu-central-2",   // Europe (Zurich)
			"eu-north-1",     // Europe (Stockholm)
			"eu-south-1",     // Europe (Milan)
			"eu-south-2",     // Europe (Spain)
			"me-south-1",     // Middle East (Bahrain)
			"me-central-1",   // Middle East (UAE)
			"il-central-1",   // Middle East (Israel)
			"af-south-1",     // Africa (Cape Town)
			"ap-southeast-1", // Asia Pacific (Singapore)
			"ap-southeast-2", // Asia Pacific (Sydney)
			"ap-southeast-3", // Asia Pacific (Jakarta)
			"ap-southeast-4", // Asia Pacific (Bangkok)
			"ap-northeast-1", // Asia Pacific (Tokyo)
			"ap-northeast-2", // Asia Pacific (Seoul)
			"ap-northeast-3", // Asia Pacific (Osaka)
			"ap-south-1",     // Asia Pacific (Mumbai)
			"ap-south-2",     // Asia Pacific (Hyderabad)
			"ap-east-1",      // Asia Pacific (Hong Kong)
			"ap-southeast-5", // Asia Pacific (Melbourne)
		},
	}
	// Run the prompt //
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputAWSBucketNamePrompt() (string, error) {
	// Prompt for AWS bucket name //
	prompt := promptui.Prompt{
		Label: "Enter AWS Bucket Name",
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputAWSAccessKeyIdPrompt() (string, error) {
	// Prompt for AWS access key id //
	prompt := promptui.Prompt{
		Label: "Enter AWS Access Key ID",
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputAWSAccessKeySecretPrompt() (string, error) {
	// Prompt for AWS access key secret //
	prompt := promptui.Prompt{
		Label: "Enter AWS Access Key Secret",
		Mask:  '*',
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func DisplayInputAWSEndpointPrompt() (string, error) {
	// Prompt for AWS endpoint //
	prompt := promptui.Prompt{
		Label:   "Enter AWS Endpoint",
		Default: "default",
	}
	// Run the prompt //
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
