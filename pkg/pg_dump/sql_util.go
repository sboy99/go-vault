package pgdump

import (
	"fmt"
	"strings"
)

func makeSqlComment(comment string) string {
	return fmt.Sprintf("--\n-- %s\n--\n", strings.Replace(comment, "\n", "\n-- ", -1))
}

func getListSchemasQuery() string {
	return `
SELECT schema_name
FROM information_schema.schemata
WHERE schema_name NOT IN ('information_schema', 'pg_catalog')
ORDER BY schema_name;
`
}

func getListTablesQuery(schema string) string {
	return fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'`, schema)
}

func getListExtensionsQuery() string {
	return `SELECT extname FROM pg_extension;`
}

func getCreateSequenceQuery(schema string) string {
	return fmt.Sprintf(`
SELECT 'CREATE SEQUENCE ' || n.nspname || '.' || c.relname || ';' as seq_creation,
    pg_get_serial_sequence(quote_ident(n.nspname) || '.' || quote_ident(t.relname), quote_ident(a.attname)) as seq_owned,
    'ALTER TABLE ' || quote_ident(n.nspname) || '.' || quote_ident(t.relname) ||
    ' ALTER COLUMN ' || quote_ident(a.attname) ||
    ' SET DEFAULT nextval(''' || n.nspname || '.' || c.relname || '''::regclass);' as col_default
FROM pg_class c
JOIN pg_namespace n ON c.relnamespace = n.oid
JOIN pg_depend d ON d.objid = c.oid AND d.deptype = 'a' AND d.classid = 'pg_class'::regclass
JOIN pg_attrdef ad ON ad.adrelid = d.refobjid AND ad.adnum = d.refobjsubid
JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
JOIN pg_class t ON t.oid = d.refobjid AND t.relkind = 'r'
WHERE c.relkind = 'S' AND t.relname = $1 AND n.nspname = '%s';
`, schema)
}

func getCreatePrimaryKeyQuery(schema string) string {
	return fmt.Sprintf(`
SELECT con.conname AS constraint_name,
    pg_get_constraintdef(con.oid) AS constraint_def
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = connamespace
WHERE con.contype = 'p'
AND rel.relname = $1
AND nsp.nspname = '%s';
`, schema)
}
