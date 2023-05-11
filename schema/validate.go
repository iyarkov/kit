package schema

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type Schema struct {
	Name         string
	Tables       map[string]Table
	Sequences    []string
	sequencesMap map[string]bool
}

type Table struct {
	Columns     map[string]Column
	Indexes     map[string]Index
	ForeignKeys map[string]ForeignKey

	name string
}

type Column struct {
	Type         string
	CharLength   int32
	NumPrecision int32
	NotNull      bool
	IsUnique     bool

	name string
}

type Index struct {
	Columns  []string
	IsUnique bool

	name       string
	columnsMap map[string]bool
}

type ForeignKey struct {
	ForeignTable string
	Columns      map[string]string

	name string
}

func Validate(ctx context.Context, db *sql.DB, expected Schema, strict bool) ([]string, error) {
	normalize(&expected)
	actual := Schema{
		Name: expected.Name,
	}
	log := zerolog.Ctx(ctx)
	log.Info().Msgf("Validating DB schema")

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// Load schema
	if err := loadTables(ctx, db, &actual); err != nil {
		return nil, fmt.Errorf("failed to load tables: %w", err)
	}
	if err := loadSequences(ctx, db, &actual); err != nil {
		return nil, fmt.Errorf("failed to load sequences: %w", err)
	}
	if err := loadColumns(ctx, db, &actual); err != nil {
		return nil, fmt.Errorf("failed to load columns: %w", err)
	}
	if err := loadIndexes(ctx, db, &actual); err != nil {
		return nil, fmt.Errorf("failed to load indexes: %w", err)
	}
	if err := loadForeignKeys(ctx, db, &actual); err != nil {
		return nil, fmt.Errorf("failed to load foreign keys: %w", err)
	}
	return validateSchema(expected, actual, strict), nil
}

func normalize(schema *Schema) {
	if schema.Tables != nil {
		for tblName, table := range schema.Tables {
			table.name = tblName
			schema.Tables[tblName] = table
			for colName, column := range table.Columns {
				column.name = colName
				table.Columns[colName] = column
			}
			for idxName, index := range table.Indexes {
				index.name = idxName
				index.columnsMap = make(map[string]bool)
				for _, colName := range index.Columns {
					index.columnsMap[colName] = true
				}
				table.Indexes[idxName] = index
			}
			for fkName, fk := range table.ForeignKeys {
				fk.name = fkName
				table.ForeignKeys[fkName] = fk
			}
		}
	}
	if schema.Sequences != nil {
		schema.sequencesMap = make(map[string]bool)
		for _, name := range schema.Sequences {
			schema.sequencesMap[name] = true
		}
	}
}

func validateSchema(expected Schema, actual Schema, strict bool) []string {
	result := validateSequences(expected, actual, strict)
	result = append(result, validateTables(expected, actual, strict)...)
	return result
}
func validateSequences(expected Schema, actual Schema, strict bool) []string {
	// Validate sequences
	result := make([]string, 0)
	for sequence := range expected.sequencesMap {
		if _, ok := actual.sequencesMap[sequence]; !ok {
			result = append(result, fmt.Sprintf("sequence %s is missing", sequence))
		}
	}
	if strict {
		for sequence := range actual.sequencesMap {
			if _, ok := expected.sequencesMap[sequence]; !ok {
				result = append(result, fmt.Sprintf("Unexpected sequence: %s", sequence))
			}
		}
	}
	return result
}

func validateTables(expected Schema, actual Schema, strict bool) []string {
	result := make([]string, 0)
	for name, expectedTable := range expected.Tables {
		if actualTable, ok := actual.Tables[name]; ok {
			result = append(result, validateTable(expectedTable, actualTable, strict)...)
		} else {
			result = append(result, fmt.Sprintf("table %s is missing", name))
		}
	}
	if strict {
		for name := range actual.Tables {
			if _, ok := expected.Tables[name]; !ok {
				result = append(result, fmt.Sprintf("Unexpected table: %s", name))
			}
		}
	}
	return result
}

func validateTable(expectedTable Table, actualTable Table, strict bool) []string {
	result := make([]string, 0)
	for name, expectedColumn := range expectedTable.Columns {
		if actualColumn, ok := actualTable.Columns[name]; ok {
			result = append(result, validateColumn(expectedTable.name, expectedColumn, actualColumn, strict)...)
		} else {
			result = append(result, fmt.Sprintf("column %s.%s is missing", expectedTable.name, name))
		}
	}
	if strict {
		for name := range actualTable.Columns {
			if _, ok := expectedTable.Columns[name]; !ok {
				result = append(result, fmt.Sprintf("Unexpected column: %s.%s", expectedTable.name, name))
			}
		}
	}

	if strict {
		for name, expectedIndex := range expectedTable.Indexes {
			if actualIndex, ok := actualTable.Indexes[name]; ok {
				result = append(result, validateIndex(expectedTable.name, expectedIndex, actualIndex)...)
			} else {
				result = append(result, fmt.Sprintf("index %s.%s is missing", expectedTable.name, name))
			}
		}
		for name := range actualTable.Indexes {
			if _, ok := expectedTable.Indexes[name]; !ok {
				result = append(result, fmt.Sprintf("Unexpected index: %s.%s", expectedTable.name, name))
			}
		}
	}

	if strict {
		for name, expectedFK := range expectedTable.ForeignKeys {
			if actualFK, ok := actualTable.ForeignKeys[name]; ok {
				result = append(result, validateForeignKey(expectedTable.name, expectedFK, actualFK)...)
			} else {
				result = append(result, fmt.Sprintf("foreign keys %s.%s is missing", expectedTable.name, name))
			}
		}
		for name := range actualTable.ForeignKeys {
			if _, ok := expectedTable.ForeignKeys[name]; !ok {
				result = append(result, fmt.Sprintf("Unexpected foreign keys: %s.%s", expectedTable.name, name))
			}
		}
	}

	return result
}

func validateColumn(tableName string, expectedColumn Column, actualColumn Column, strict bool) []string {
	result := make([]string, 0)
	if expectedColumn.Type != actualColumn.Type {
		result = append(result, fmt.Sprintf("invalid column type: %s.%s, expected %s, actual %s", tableName, expectedColumn.name, expectedColumn.Type, actualColumn.Type))
	}
	if strict {
		if expectedColumn.CharLength != actualColumn.CharLength {
			result = append(result, fmt.Sprintf("invalid column char length: %s.%s, expected %d, actual %d", tableName, expectedColumn.name, expectedColumn.CharLength, actualColumn.CharLength))
		}
		if expectedColumn.NumPrecision != actualColumn.NumPrecision {
			result = append(result, fmt.Sprintf("invalid column num precision: %s.%s, expected %d, actual %d", tableName, expectedColumn.name, expectedColumn.NumPrecision, actualColumn.NumPrecision))
		}
		if expectedColumn.NotNull != actualColumn.NotNull {
			result = append(result, fmt.Sprintf("invalid column is nullable: %s.%s, expected %t, actual %t", tableName, expectedColumn.name, expectedColumn.NotNull, actualColumn.NotNull))
		}
		if expectedColumn.IsUnique != actualColumn.IsUnique {
			result = append(result, fmt.Sprintf("invalid column is unique: %s.%s, expected %t, actual %t", tableName, expectedColumn.name, expectedColumn.IsUnique, actualColumn.IsUnique))
		}
	}
	return result
}

func validateIndex(tableName string, expectedIndex Index, actualIndex Index) []string {
	result := make([]string, 0)
	if expectedIndex.IsUnique != actualIndex.IsUnique {
		result = append(result, fmt.Sprintf("invalid index IsUnique: %s.%s, expected %t, actual %t", tableName, expectedIndex.name, expectedIndex.IsUnique, actualIndex.IsUnique))
	}
	for expectedColumn := range expectedIndex.columnsMap {
		if _, ok := actualIndex.columnsMap[expectedColumn]; !ok {
			result = append(result, fmt.Sprintf("invalid index  %s.%s, missing column: %s", tableName, expectedIndex.name, expectedColumn))
		}
	}
	for actualColumn := range actualIndex.columnsMap {
		if _, ok := expectedIndex.columnsMap[actualColumn]; !ok {
			result = append(result, fmt.Sprintf("invalid index  %s.%s, extra column: %s", tableName, expectedIndex.name, actualColumn))
		}
	}
	return result
}

func validateForeignKey(tableName string, expectedFK ForeignKey, actualFK ForeignKey) []string {
	result := make([]string, 0)
	if expectedFK.ForeignTable != actualFK.ForeignTable {
		result = append(result, fmt.Sprintf("invalid fk foreign table: %s.%s, expected %s, actual %s", tableName, expectedFK.name, expectedFK.ForeignTable, actualFK.ForeignTable))
	}
	for column, expectedForeignColumn := range expectedFK.Columns {
		if actualForeignColumn, ok := actualFK.Columns[column]; !ok {
			result = append(result, fmt.Sprintf("invalid fk: %s.%s, missed column: %s => %s", tableName, expectedFK.name, column, expectedForeignColumn))
		} else if actualForeignColumn != expectedForeignColumn {
			result = append(result, fmt.Sprintf("invalid fk: %s.%s, wrong column mapping, expected: %s => %s, actual %s => %s", tableName, expectedFK.name, column, expectedForeignColumn, column, actualForeignColumn))
		}
	}
	for column, actualForeignColumn := range actualFK.Columns {
		if _, ok := expectedFK.Columns[column]; !ok {
			result = append(result, fmt.Sprintf("invalid fk: %s.%s, extra column %s => %s", tableName, expectedFK.name, column, actualForeignColumn))
		}
	}
	return result
}
