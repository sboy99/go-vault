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
	tables, err := p.listTables()
	if err != nil {
		return nil, err
	}

	var dumpSql bytes.Buffer
	for _, table := range tables {
		dumpSqlStatement, err := p.getDumpSqlStatement(table)
		if err != nil {
			return nil, err
		}
		dumpSql.WriteString(dumpSqlStatement)
		dumpSql.WriteString("\n\n")
	}

	return dumpSql.Bytes(), nil
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

func (p *PgDump) getDumpSqlStatement(table string) (string, error) {
	createTableStatement, err := p.getCreateTableStatement(table)
	if err != nil {
		return "", err
	}

	return createTableStatement, nil
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
