package schema

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"microserver.rockyrunstream.com/foundation/support"
)

func loadTables(ctx context.Context, db *sql.DB, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = db.QueryContext(ctx, queryLoadTable, schema.Name)
		defer support.CloseWithWarning(ctx, rows, "Failed to close load tables query")
		if err != nil {
			return fmt.Errorf("failed to query tables: %w", err)
		}
	}

	var name string
	tables := make(map[string]Table)
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan rows: %w", err)
		}
		tables[name] = Table{
			name:        name,
			Columns:     make(map[string]Column),
			Indexes:     make(map[string]Index),
			ForeignKeys: make(map[string]ForeignKey),
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("after scan error: %w", rows.Err())
	}

	schema.Tables = tables
	return nil
}

func loadSequences(ctx context.Context, db *sql.DB, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = db.QueryContext(ctx, queryLoadSequences, schema.Name)
		defer support.CloseWithWarning(ctx, rows, "Failed to close load sequences query")
		if err != nil {
			return fmt.Errorf("failed to query sequences: %w", err)
		}
	}

	var name string
	sequences := make(map[string]bool)
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		sequences[name] = true
	}

	if rows.Err() != nil {
		return fmt.Errorf("after scan error: %w", rows.Err())
	}

	schema.sequencesMap = sequences

	i := 0
	schema.Sequences = make([]string, len(schema.sequencesMap))
	for sequenceName, _ := range sequences {
		schema.Sequences[i] = sequenceName
		i++
	}
	return nil
}

func loadColumns(ctx context.Context, db *sql.DB, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = db.QueryContext(ctx, queryLoadColumns, schema.Name)
		defer support.CloseWithWarning(ctx, rows, "Failed to close load columns query")
		if err != nil {
			return fmt.Errorf("failed to query columns: %w", err)
		}
	}

	var column Column
	var tableName string
	var charLength sql.NullInt32
	var numPrecision sql.NullInt32
	var isNullable bool

	log := zerolog.Ctx(ctx)
	for rows.Next() {
		if err := rows.Scan(&tableName, &column.name, &column.Type, &charLength, &numPrecision, &isNullable, &column.IsUnique); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		if charLength.Valid {
			column.CharLength = charLength.Int32
		} else {
			column.CharLength = 0
		}
		if numPrecision.Valid {
			column.NumPrecision = numPrecision.Int32
		} else {
			column.NumPrecision = 0
		}
		column.NotNull = !isNullable
		table, ok := schema.Tables[tableName]
		if ok {
			table.Columns[column.name] = column
		} else {
			log.Error().Msgf("table %s not found for column %s", tableName, column.name)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("after scan error: %w", rows.Err())
	}

	return nil
}

func loadIndexes(ctx context.Context, db *sql.DB, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = db.QueryContext(ctx, queryLoadIndexes, schema.Name)
		defer support.CloseWithWarning(ctx, rows, "Failed to close load indexes query")
		if err != nil {
			return fmt.Errorf("failed to query indexes: %w", err)
		}
	}
	var tableName, indexName, columnName string
	var isUnique bool

	log := zerolog.Ctx(ctx)
	for rows.Next() {
		if err := rows.Scan(&tableName, &indexName, &columnName, &isUnique); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		table, ok := schema.Tables[tableName]
		if ok {
			index := table.Indexes[indexName]
			index.name = indexName
			index.IsUnique = isUnique
			if index.columnsMap == nil {
				index.columnsMap = make(map[string]bool)
				index.Columns = make([]string, 0)
			}
			index.columnsMap[columnName] = true
			index.Columns = append(index.Columns, columnName)
			table.Indexes[indexName] = index
		} else {
			log.Error().Msgf("table %s not found for index %s", tableName, indexName)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("after scan error: %w", rows.Err())
	}

	return nil
}

func loadForeignKeys(ctx context.Context, db *sql.DB, schema *Schema) error {
	var rows *sql.Rows
	{
		var err error
		rows, err = db.QueryContext(ctx, queryLoadForeignKeys, schema.Name)
		defer support.CloseWithWarning(ctx, rows, "Failed to close load foreign keys query")
		if err != nil {
			return fmt.Errorf("failed to query foreign keys: %w", err)
		}
	}

	var tableName, keyName, columnName, foreignTableName, foreignColumnName string

	log := zerolog.Ctx(ctx)
	for rows.Next() {
		if err := rows.Scan(&tableName, &keyName, &columnName, &foreignTableName, &foreignColumnName); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		table, ok := schema.Tables[tableName]
		if ok {
			fk := table.ForeignKeys[tableName]
			fk.name = keyName
			fk.ForeignTable = foreignTableName
			if fk.Columns == nil {
				fk.Columns = make(map[string]string)
			}
			fk.Columns[columnName] = foreignColumnName
		} else {
			log.Error().Msgf("table %s not found for foreign key %s", tableName, keyName)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("after scan error: %w", rows.Err())
	}
	return nil
}
