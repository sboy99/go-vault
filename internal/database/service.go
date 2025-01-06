package database

import (
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/storage"
	"github.com/sboy99/go-vault/pkg/logger"
)

type DatabaseService struct {
	storage  *storage.Storage
	database *Database
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		storage:  storage.NewStorage(),
		database: NewDatabase(),
	}
}

func (d *DatabaseService) Backup() error {
	cfg := config.GetConfig()

	logger.Info("Connecting to DB...")
	if err := d.database.Connect(
		cfg.DB.Type,
		cfg.DB.Name,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
	); err != nil {
		return err
	}
	defer d.database.Close(cfg.DB.Type)
	logger.Info("Connected to DB")

	logger.Info("Pinging DB...")
	if err := d.database.Ping(cfg.DB.Type); err != nil {
		return err
	}
	logger.Info("Pinged DB")

	logger.Info("Backing up DB...")
	data, err := d.database.Backup(cfg.DB.Type)
	if err != nil {
		return err
	}
	logger.Info("Backed up DB")

	logger.Info("Saving backup...")
	if err := d.storage.Save(cfg.Storage.Type, "backup.sql", data); err != nil {
		return err
	}
	logger.Info("Saved backup")

	return nil
}
