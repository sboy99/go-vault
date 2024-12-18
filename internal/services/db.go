package services

import (
	"github.com/sboy99/go-vault/internal/strategies"
	"github.com/sboy99/go-vault/pkg/logger"
)

type DatabaseService struct {
	storage  *strategies.Storage
	database *strategies.Database
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		storage:  strategies.NewStorage("."),
		database: strategies.NewDatabase(),
	}
}

func (d *DatabaseService) Backup(dbType strategies.DatabaseEnum, name string, host string, port int, username string, password string) error {
	logger.Info("Connecting to DB...")
	if err := d.database.Connect(dbType, name, host, port, username, password); err != nil {
		return err
	}
	defer d.database.Close(dbType)
	logger.Info("Connected to DB")

	logger.Info("Pinging DB...")
	if err := d.database.Ping(dbType); err != nil {
		return err
	}
	logger.Info("Pinged DB")

	logger.Info("Backing up DB...")
	data, err := d.database.Backup(dbType)
	if err != nil {
		return err
	}
	logger.Info("Backed up DB")

	logger.Info("Saving backup...")
	if err := d.storage.Save(strategies.FILE, "backup.sql", data); err != nil {
		return err
	}
	logger.Info("Saved backup")

	return nil
}
