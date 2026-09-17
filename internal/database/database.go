package database

import (
	"fmt"

	"github.com/sboy99/go-vault/config"
)

// IDatabase is retained for CLI compatibility; Postgres now uses the engine via backup.Service.
type IDatabase interface {
	Connect(name string, host string, port int, username string, password string) error
	Ping() error
	Close() error
}

type Database struct {
	dbMap map[config.DatabaseEnum]IDatabase
}

func NewDatabase() *Database {
	return &Database{
		dbMap: map[config.DatabaseEnum]IDatabase{
			config.POSTGRESQL: NewPostgresDB(),
		},
	}
}

func (d *Database) getDatabase(dbType config.DatabaseEnum) (IDatabase, error) {
	db, ok := d.dbMap[dbType]
	if !ok || db == nil {
		return nil, fmt.Errorf("unsupported database type %q", dbType)
	}
	return db, nil
}

func (d *Database) Connect(dbType config.DatabaseEnum, name string, host string, port int, username string, password string) error {
	db, err := d.getDatabase(dbType)
	if err != nil {
		return err
	}
	return db.Connect(name, host, port, username, password)
}

func (d *Database) Close(dbType config.DatabaseEnum) error {
	db, err := d.getDatabase(dbType)
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *Database) Ping(dbType config.DatabaseEnum) error {
	db, err := d.getDatabase(dbType)
	if err != nil {
		return err
	}
	return db.Ping()
}
