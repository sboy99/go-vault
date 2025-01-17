package database

import (
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/internal/storage"
	"github.com/sboy99/go-vault/pkg/logger"
	"github.com/sboy99/go-vault/pkg/utils"
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

func (d *DatabaseService) Backup() {
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
		logger.Error("Failed to connect to DB\nDetails: %v", err)
		return
	}
	defer d.database.Close(cfg.DB.Type)
	logger.Info("Connected to DB")

	logger.Info("Pinging DB...")
	if err := d.database.Ping(cfg.DB.Type); err != nil {
		logger.Error("Failed to ping DB\nDetails: %v", err)
		return
	}
	logger.Info("Pinged DB")

	logger.Info("Backing up DB...")
	data, err := d.database.Backup(cfg.DB.Type)
	if err != nil {
		logger.Error("Failed to backup DB\nDetails: %v", err)
		return
	}
	logger.Info("Backed up DB")

	logger.Info("Saving backup...")
	backupFilename := buidlFileName(cfg.DB.Name)
	if err := d.storage.Save(cfg.Storage.Type, backupFilename, data); err != nil {
		logger.Error("Failed to save backup\nDetails: %v", err)
		return
	}
	logger.Info("Saved backup")

	backupMeta := meta.NewBackupMeta(backupFilename, cfg.DB.Type, cfg.Storage.Type)
	if err := backupMeta.Save(); err != nil {
		logger.Error("Failed to save backup meta\nDetails: %v", err)
		return
	}
	logger.Info("Backup successful")
}

func buidlFileName(dbName string) string {
	return utils.GetUnixTimeStamp() + "_" + dbName + "_" + "backup" + ".sql"
}
