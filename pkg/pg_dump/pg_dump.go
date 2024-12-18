package pgdump

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type PgDump struct {
	db *sql.DB
}

func NewPgDump(db *sql.DB) *PgDump {
	return &PgDump{
		db: db,
	}
}

func (p *PgDump) Dump() ([]byte, error) {
	dbVersion, err := p.getDatabaseVersion()
	if err != nil {
		return nil, err
	}

	tables, err := p.listTables()
	if err != nil {
		return nil, err
	}

	var dumpSql bytes.Buffer
	dumpSql.WriteString(makeSqlComment("PostgreSQL database dump"))
	dumpSql.WriteString(makeSqlComment(fmt.Sprintf("Database version: %s", dbVersion)))
	dumpSql.WriteString("SET statement_timeout = 0;\n")
	dumpSql.WriteString("SET lock_timeout = 0;\n")
	dumpSql.WriteString("SET client_encoding = 'UTF8';\n")
	dumpSql.WriteString("SET standard_conforming_strings = on;\n")
	dumpSql.WriteString("SET check_function_bodies = FALSE;\n")
	dumpSql.WriteString("SET client_min_messages = warning;\n\n")

	createTableStatement, err := p.generateCreateTableStatement(tables)
	if err != nil {
		return nil, err
	}
	dumpSql.WriteString(createTableStatement)

	return dumpSql.Bytes(), nil
}

func (p *PgDump) getDatabaseVersion() (string, error) {
	var version string
	if err := p.db.QueryRow("SELECT version()").Scan(&version); err != nil {
		return "", err
	}
	return version, nil
}

func (p *PgDump) listTables() ([]string, error) {
	rows, err := p.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (p *PgDump) generateCreateTableStatement(tables []string) (string, error) {
	var dumpSql bytes.Buffer
	for _, table := range tables {
		createTableStatement, err := p.getCreateTableStatement(table)
		if err != nil {
			return "", err
		}
		dumpSql.WriteString(createTableStatement)
		dumpSql.WriteString("\n\n")
	}
	return dumpSql.String(), nil
}

func (p *PgDump) getCreateTableStatement(tableName string) (string, error) {
	query := fmt.Sprintf("SELECT column_name, data_type, character_maximum_length FROM information_schema.columns WHERE table_name = '%s'", tableName)
	rows, err := p.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName, dataType string
		var charMaxLength *int
		if err := rows.Scan(&columnName, &dataType, &charMaxLength); err != nil {
			return "", err
		}
		columnDef := fmt.Sprintf("%s %s", columnName, dataType)
		if charMaxLength != nil {
			columnDef += fmt.Sprintf("(%d)", *charMaxLength)
		}
		columns = append(columns, columnDef)
	}

	return fmt.Sprintf("CREATE TABLE %s (\n    %s\n);", tableName, strings.Join(columns, ",\n    ")), nil
}
