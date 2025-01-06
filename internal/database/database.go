package database

import "github.com/sboy99/go-vault/config"

type IDatabase interface {
	Connect(name string, host string, port int, username string, password string) error
	Ping() error
	Backup() ([]byte, error)
	Restore(data []byte) error
	Close() error
}

type Database struct {
	dbMap map[config.DatabaseEnum]IDatabase
}

func NewDatabase() *Database {
	return &Database{
		dbMap: map[config.DatabaseEnum]IDatabase{
			config.POSTGRESQL: NewPostgresDB(),
			config.MYSQL:      nil,
			config.MONGODB:    nil,
		},
	}
}

func (d *Database) Connect(dbType config.DatabaseEnum, name string, host string, port int, username string, password string) error {
	db := d.getDatabase(dbType)
	if err := db.Connect(name, host, port, username, password); err != nil {
		return err
	}
	return nil
}

func (d *Database) Close(dbType config.DatabaseEnum) error {
	db := d.getDatabase(dbType)
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Ping(dbType config.DatabaseEnum) error {
	db := d.getDatabase(dbType)
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Backup(dbType config.DatabaseEnum) ([]byte, error) {
	db := d.getDatabase(dbType)
	data, err := db.Backup()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Database) Restore(dbType config.DatabaseEnum, data []byte) error {
	db := d.getDatabase(dbType)
	if err := db.Restore(data); err != nil {
		return err
	}
	return nil
}

func (d *Database) getDatabase(dbType config.DatabaseEnum) IDatabase {
	return d.dbMap[dbType]
}
