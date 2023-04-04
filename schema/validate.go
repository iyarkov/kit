package schema

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type Schema struct {
	oid       int64
	Name      string
	Tables    []Table
	Sequences []string
}

type Table struct {
	oid         int64
	Name        string
	Columns     []Column
	Indexes     []Index
	ForeignKeys []ForeignKey
}

type Column struct {
	oid  int64
	Name string
	Type string
}

type Index struct {
	oid     int64
	Name    string
	Columns []string
}

type ForeignKey struct {
	oid           int64
	Name          string
	Column        string
	ForeignColumn string
}

func Validate(ctx context.Context, db *sql.DB, expected Schema) ([]string, error) {
	actual := Schema{
		Name: expected.Name,
	}
	log := zerolog.Ctx(ctx)
	log.Info().Msgf("Validating DB schema")

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// Start transaction
	var tx *sql.Tx
	{
		var err error
		tx, err = db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to begint transaction: %w", err)
		}
		defer func() {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				zerolog.Ctx(ctx).Error().Err(rollbackErr).Msg("tx rollback error")
			}
		}()
	}

	// Start tables
	if err := loadTables(ctx, tx, &actual); err != nil {
		return nil, fmt.Errorf("failed to load tables: %w", err)
	}

	return nil, fmt.Errorf("not implemented")
}

var loadTableQuery = `select table_name from information_schema."tables" t where table_schema = ?`

func loadTables(ctx context.Context, tx *sql.Tx, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = tx.QueryContext(ctx, loadTableQuery, schema.Name)
		if err != nil {
			return fmt.Errorf("failed to query tables: %w", err)
		}
	}

	var name string
	tables := make([]Table, 0)
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan result: %w", err)
		}
		tables = append(tables, Table{Name: name})
	}

	schema.Tables = tables
	return nil
}
