package strategies

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/lib/pq"
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
	var cmdout, cmderr bytes.Buffer
	cmd := exec.Command("pg_dump", "--host", p.host, "--port", fmt.Sprintf("%d", p.port), "--username", p.username, "--dbname", p.name)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", p.password))
	cmd.Stdout = &cmdout
	cmd.Stderr = &cmderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return cmdout.Bytes(), nil
}

func (p *PostgresDB) Restore(data []byte) error {
	return nil
}
