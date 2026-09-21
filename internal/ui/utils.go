package ui

import "github.com/sboy99/go-vault/internal/config"

func getDatabasePort(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "5432"
	default:
		return ""
	}
}

func getDatabaseUser(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "postgres"
	default:
		return ""
	}
}

func getDatabaseName(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "postgres"
	default:
		return ""
	}
}
