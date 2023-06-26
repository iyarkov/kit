package sql

// Schema History queries
var queryTableExists = `
SELECT EXISTS (
	SELECT FROM pg_tables WHERE
schemaname = $1 AND
tablename  = 'schema_history'
)`

var queryCreateTable = `
CREATE TABLE schema_history (
	id SERIAL,
	created_at TIMESTAMP(3) WITHOUT TIME ZONE,
	version VARCHAR(255),
	PRIMARY KEY (id)
)`

var queryLastVersion = "SELECT version FROM schema_history WHERE id = (SELECT MAX(id) FROM schema_history)"

var queryInsertVersion = "INSERT INTO schema_history(created_at, version) VALUES($1, $2)"

var queryLoadHistory = "SELECT id, created_at, version FROM schema_history ORDER BY id DESC"

// Schema Validation queries

var queryLoadTable = `SELECT tablename FROM pg_tables WHERE schemaname = $1  and tablename != 'schema_history'`

var queryLoadColumns = `SELECT c.table_name, c.column_name, c.udt_name, c.character_maximum_length, c.numeric_precision,
       CASE
           WHEN c.is_nullable = 'YES' THEN true
           WHEN c.is_nullable = 'NO' THEN false
       END AS is_nullable,
       CASE
            WHEN tc.constraint_type = 'UNIQUE' THEN true
            ELSE false
       END AS is_unique
FROM information_schema.columns c LEFT JOIN
    information_schema.key_column_usage kcu ON c.table_name = kcu.table_name AND c.column_name = kcu.column_name LEFT JOIN
    information_schema.table_constraints tc ON kcu.constraint_name = tc.constraint_name
WHERE c.table_schema = $1 and c.table_name != 'schema_history'
ORDER BY c.table_name, c.ordinal_position`

var queryLoadSequences = `SELECT sequencename FROM pg_sequences where schemaname=$1 and sequencename != 'schema_history_id_seq' order by sequencename`

var queryLoadIndexes = `SELECT
    t.relname AS table_name,
    i.relname AS index_name,
    a.attname AS column_name,
    ix.indisunique as is_unique
FROM
    pg_class t
        JOIN
    pg_index ix ON t.oid = ix.indrelid
        JOIN
    pg_class i ON ix.indexrelid = i.oid
        JOIN
    pg_attribute a ON t.oid = a.attrelid AND a.attnum = ANY(ix.indkey)
        JOIN
    pg_namespace n ON t.relnamespace = n.oid
WHERE
        t.relkind = 'r' -- Only relational tables (excluding materialized views and other types)
        and n.nspname = $1
        and i.relname != 'schema_history_pkey'
ORDER BY
    n.nspname, t.relname, i.relname`

var queryLoadForeignKeys = `SELECT
    tc.table_name AS table_name,
    tc.constraint_name AS key_name,
    kcu.column_name AS column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM
    information_schema.table_constraints tc
        JOIN
    information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
        JOIN
    information_schema.constraint_column_usage ccu ON tc.constraint_name = ccu.constraint_name
WHERE
        tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = $1`

var queryDropSchema = "DROP SCHEMA %s CASCADE"
var queryCreateSchema = "CREATE SCHEMA %s"
