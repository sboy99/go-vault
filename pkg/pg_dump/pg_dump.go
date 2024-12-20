package pgdump

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sboy99/go-vault/pkg/logger"
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

	createTableStatements, err := p.generateCreateTableStatementsForTables(tables)
	if err != nil {
		return nil, err
	}
	dumpSql.WriteString(createTableStatements)

	createSequenceStatements, err := p.generateCreateSequenceStatenentsForTables(tables)
	if err != nil {
		return nil, err
	}
	dumpSql.WriteString(createSequenceStatements)

	createPrimaryKeyStatements, err := p.generateCreatePrimaryKeyStatementsForTables(tables)
	if err != nil {
		return nil, err
	}
	dumpSql.WriteString(createPrimaryKeyStatements)

	tableDataCopyStatements, err := p.copyDataOfTables(tables)
	if err != nil {
		return nil, err
	}
	dumpSql.WriteString(tableDataCopyStatements)

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

func (p *PgDump) generateCreateTableStatementsForTables(tables []string) (string, error) {
	var dumpSql bytes.Buffer
	dumpSql.WriteString(makeSqlComment("Create Tables"))
	for _, table := range tables {
		createTableStatement, err := p.getCreateTableStatement(table)
		if err != nil {
			return "", fmt.Errorf("error getting create table statement for table %s: %v", table, err)
		}
		if createTableStatement == "" {
			continue
		}
		dumpSql.WriteString(createTableStatement)
		dumpSql.WriteString("\n\n")
	}
	return dumpSql.String(), nil
}

func (p *PgDump) generateCreatePrimaryKeyStatementsForTables(tables []string) (string, error) {
	var dumpSql bytes.Buffer
	dumpSql.WriteString(makeSqlComment("Create Primary Keys"))
	for _, table := range tables {
		pkStatement, err := p.getCreatePrimaryKeyStatement(table)
		if err != nil {
			return "", fmt.Errorf("error getting create primary key statement for table %s: %v", table, err)
		}
		if pkStatement == "" {
			continue
		}
		dumpSql.WriteString(pkStatement)
		dumpSql.WriteString("\n\n")
	}
	return dumpSql.String(), nil
}

func (p *PgDump) generateCreateSequenceStatenentsForTables(tables []string) (string, error) {
	var dumpSql bytes.Buffer
	dumpSql.WriteString(makeSqlComment("Create Sequences"))
	for _, table := range tables {
		sequences, err := p.getCreateSequenceStatement(table)
		if err != nil {
			return "", err
		}
		if sequences == "" {
			continue
		}
		dumpSql.WriteString(sequences)
		dumpSql.WriteString("\n\n")
	}
	return dumpSql.String(), nil
}

func (p *PgDump) copyDataOfTables(tables []string) (string, error) {
	var dumpSql bytes.Buffer
	dumpSql.WriteString(makeSqlComment("Copy Table Data"))
	for _, table := range tables {
		tableDataCopyStatement, err := p.getTableDataCopyStatement(table)
		if err != nil {
			return "", fmt.Errorf("error getting table data copy statement for table %s: %v", table, err)
		}
		if tableDataCopyStatement == "" {
			continue
		}
		dumpSql.WriteString(tableDataCopyStatement)
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

func (p *PgDump) getCreateSequenceStatement(tableName string) (string, error) {
	var sequencesSQL strings.Builder

	rows, err := p.db.Query(getCreateSequenceQuery("public"), tableName)
	if err != nil {
		return "", fmt.Errorf("error querying sequences for table %s: %v", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var seqCreation, seqOwned, colDefault string
		if err := rows.Scan(&seqCreation, &seqOwned, &colDefault); err != nil {
			return "", fmt.Errorf("error scanning sequence information: %v", err)
		}

		// Here we directly use the sequence creation script.
		// The seqOwned might not be necessary if we're focusing on creation and default value setting.
		sequencesSQL.WriteString(seqCreation + "\n" + colDefault + "\n")
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating over sequences: %v", err)
	}

	return sequencesSQL.String(), nil
}

func (p *PgDump) getCreatePrimaryKeyStatement(tableName string) (string, error) {
	var pksSQL strings.Builder
	rows, err := p.db.Query(getCreatePrimaryKeyQuery("public"), tableName)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var constraintName, constraintDef string
		if err := rows.Scan(&constraintName, &constraintDef); err != nil {
			return "", fmt.Errorf("error scanning primary key information: %v", err)
		}
		pksSQL.WriteString(fmt.Sprintf("ALTER TABLE public.%s ADD CONSTRAINT %s %s;\n",
			tableName, constraintName, constraintDef))
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating over primary keys: %v", err)
	}

	return pksSQL.String(), nil
}

func (p *PgDump) getTableDataCopyStatement(tableName string) (string, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := p.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("COPY %s (%s) FROM stdin;\n", tableName, strings.Join(columns, ", ")))
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}
		var valueStrings []string
		for _, value := range values {
			valueStrings = append(valueStrings, string(value))
		}
		logger.Debug("Row: %v", valueStrings)
		output.WriteString(strings.Join(valueStrings, "\t") + "\n")
	}
	output.WriteString("\\.\n")

	return output.String(), nil
}
