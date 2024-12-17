package services

import (
	"github.com/sboy99/go-vault/internal/strategies"
	"github.com/sboy99/go-vault/pkg/logger"
)

type DatabaseService struct {
	Storage  *strategies.Storage
	Database *strategies.Database
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		Storage:  strategies.NewStorage("."),
		Database: strategies.NewDatabase(),
	}
}

func (d *DatabaseService) Connect(dbType strategies.DatabaseEnum, name string, host string, port int, username string, password string) error {
	logger.Info("Connecting to DB...")
	if err := d.Database.Connect(dbType, name, host, port, username, password); err != nil {
		return err
	}
	logger.Info("Connected to DB")
	logger.Info("Pinging DB...")
	if err := d.Database.Ping(dbType); err != nil {
		return err
	}
	logger.Info("Pinged DB")
	logger.Info("Backing up DB...")
	data, err := d.Database.Backup(dbType)
	if err != nil {
		return err
	}
	logger.Info("Backed up DB")
	logger.Info("Saving backup...")
	if err := d.Storage.Save(strategies.FILE, "backup.sql", data); err != nil {
		return err
	}
	logger.Info("Saved backup")
	return nil
}
