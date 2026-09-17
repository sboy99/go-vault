package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sboy99/go-vault/config"
)

// PostgresDB provides a lightweight connectivity check used by health probes.
type PostgresDB struct {
	name     string
	host     string
	port     int
	username string
	password string
	sslmode  string
	db       *sql.DB
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (p *PostgresDB) Connect(name string, host string, port int, username string, password string) error {
	cfg := config.GetConfig()
	p.name = name
	p.host = host
	p.port = port
	p.username = username
	p.password = password
	p.sslmode = cfg.DB.SSLMode
	if p.sslmode == "" {
		p.sslmode = "require"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, name, p.sslmode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgresDB) Close() error {
	if p.db == nil {
		return nil
	}
	return p.db.Close()
}

func (p *PostgresDB) Ping() error {
	if p.db == nil {
		return fmt.Errorf("not connected")
	}
	return p.db.Ping()
}
