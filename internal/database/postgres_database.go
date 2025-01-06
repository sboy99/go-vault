package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	pgdump "github.com/sboy99/go-vault/pkg/pg_dump"
)

type PostgresDB struct {
	name     string
	host     string
	port     int
	username string
	password string
	db       *sql.DB
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (p *PostgresDB) Connect(name string, host string, port int, username string, password string) error {
	p.name = name
	p.host = host
	p.port = port
	p.username = username
	p.password = password

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgresDB) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Ping() error {
	if err := p.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Backup() ([]byte, error) {
	pgDump := pgdump.NewPgDump(p.db)
	return pgDump.Dump()
}

func (p *PostgresDB) Restore(data []byte) error {
	sqlContent := string(data)
	statements := strings.Split(sqlContent, ";")
	for _, statement := range statements {
		if _, err := p.db.Exec(statement); err != nil {
			return fmt.Errorf("Error executing statement: %s\nErr: %v", statement, err)
		}
	}
	return nil
}
