package ui

import "github.com/sboy99/go-vault/config"

func getDatabasePort(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "5432"
	case config.MYSQL:
		return "3306"
	case config.MONGODB:
		return "27017"
	default:
		return ""
	}
}

func getDatabaseUser(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "postgres"
	case config.MYSQL:
		return "mysql"
	case config.MONGODB:
		return "mongo"
	default:
		return ""
	}
}

func getDatabaseName(dbType config.DatabaseEnum) string {
	switch dbType {
	case config.POSTGRESQL:
		return "postgres"
	case config.MYSQL:
		return "mysql"
	case config.MONGODB:
		return "mongo"
	default:
		return ""
	}
}
