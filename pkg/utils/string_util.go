package utils

import (
	"regexp"
	"strings"
)

type FileType string

const (
	CONTROLLER      FileType = "controller"
	SERVICE         FileType = "service"
	REPOSITORY      FileType = "repository"
	ADAPTER         FileType = "adapter"
	CONTROLLER_PORT FileType = "controller.port"
	SERVICE_PORT    FileType = "service.port"
	REPOSITORY_PORT FileType = "repository.port"
	ADAPTER_PORT    FileType = "adapter.port"
)
const _FILE_EXTENSION string = ".ts"

func ToFileName(name string, fileType FileType) string {
	switch fileType {
	case CONTROLLER_PORT, SERVICE_PORT, REPOSITORY_PORT, ADAPTER_PORT, CONTROLLER, SERVICE, REPOSITORY, ADAPTER:
		return fmtKebabCase(name) + "." + string(fileType) + _FILE_EXTENSION
	default:
		return fmtKebabCase(name) + "." + string(fileType) + _FILE_EXTENSION
	}
}

func ToClassName(name string, fileType FileType) string {
	switch fileType {
	case CONTROLLER_PORT, SERVICE_PORT, REPOSITORY_PORT, ADAPTER_PORT, CONTROLLER, SERVICE, REPOSITORY, ADAPTER:
		return fmtPascalCase(name) + fmtPascalCase(string(fileType))
	default:
		return fmtPascalCase(name)
	}
}

func fmtKebabCase(input string) string {
	// Replace all non-alphanumeric characters (except underscores) with spaces
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	input = reg.ReplaceAllString(input, " ")

	// Replace underscores with spaces
	input = strings.ReplaceAll(input, "_", " ")

	// Use regex to identify words in camelCase or PascalCase and separate them
	camelCaseReg := regexp.MustCompile(`([a-z])([A-Z])`)
	input = camelCaseReg.ReplaceAllString(input, `$1 $2`)

	// Convert to lowercase and trim spaces
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)

	// Replace spaces with hyphens to create kebab case
	return strings.ReplaceAll(input, " ", "-")
}

func fmtPascalCase(input string) string {
	// Remove common delimiters by replacing them with spaces
	input = strings.ReplaceAll(input, "-", " ")
	input = strings.ReplaceAll(input, "_", " ")
	input = strings.ReplaceAll(input, ".", " ")

	// Split the input into words by spaces
	words := strings.Fields(input)

	// Capitalize the first letter of each word and join them
	var pascalCase string
	for _, word := range words {
		if len(word) > 0 {
			// Capitalize the first letter and append the rest of the word
			pascalCase += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return pascalCase
}
