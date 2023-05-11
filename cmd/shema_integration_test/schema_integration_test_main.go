package main

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"microserver.rockyrunstream.com/foundation/logger"
	"microserver.rockyrunstream.com/foundation/schema"
	"microserver.rockyrunstream.com/foundation/support"
)

var changeset = make([]schema.Change, 0)
var expectedSchema = schema.Schema{}

func main() {
	ctx := context.Background()
	logger.InitLogger(logger.Configuration{})
	db, err := sql.Open("pgx", "host=127.0.0.1 port=5432 user=postgres password=example dbname=test sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("DB connection failed")
	}
	defer support.CloseWithWarning(ctx, db, "failed to close db")
	log.Info().Msg("connected to the db")

	ensurePostgresSQLSchema(ctx, db)
	schemaV001(ctx, db)
}

func ensurePostgresSQLSchema(ctx context.Context, db *sql.DB) {
	if _, err := db.ExecContext(ctx, "drop table if exists schema_history"); err != nil {
		log.Fatal().Err(err).Msg("drop schema_history failed")
	}
	if _, err := db.ExecContext(ctx, "drop schema if exists public cascade"); err != nil {
		log.Fatal().Err(err).Msg("drop schema failed")
	}
	if _, err := db.ExecContext(ctx, "create schema public"); err != nil {
		log.Fatal().Err(err).Msg("create schema failed")
	}
	log.Info().Msg("schema created")
}

func schemaV001(ctx context.Context, db *sql.DB) {
	changeset = append(changeset, schema.Change{
		Version: "0.0.1",
		Commands: []string{
			`CREATE TABLE table_a (
				id SERIAL,
				created_at TIMESTAMP(3) WITHOUT TIME ZONE,
				name varchar(255) CONSTRAINT name_idx UNIQUE,
				description VARCHAR(255),
				b_flag bool NOT NULL,
				partition int,
				subpartition float,
				PRIMARY KEY (id)
			)`,
			"CREATE INDEX partition_idx ON table_a(partition, subpartition)",
		},
	})

	dbVersion, expectedVersion, err := schema.Update(ctx, db, changeset)
	if err != nil {
		log.Fatal().Err(err).Msg("schemaV001 method failed with unexpected error")
	}
	if dbVersion != "" {
		log.Fatal().Err(err).Msgf("dbVersion expected to be an empty string, actual value is %s", dbVersion)
	}
	if expectedVersion != "0.0.1" {
		log.Fatal().Err(err).Msgf("invalid expectedVersion expected: 0.0.1, actual [%s]", expectedVersion)
	}

	expectedSchema.Name = "public"
	expectedSchema.Sequences = []string{
		"table_a_id_seq",
	}
	expectedSchema.Tables = make(map[string]schema.Table)
	expectedSchema.Tables["table_a"] = schema.Table{
		Columns: map[string]schema.Column{
			"id": {
				Type:         "int4",
				NumPrecision: 32,
				NotNull:      true,
			},
			"created_at": {
				Type: "timestamp",
			},
			"name": {
				Type:       "varchar",
				CharLength: 255,
				IsUnique:   true,
			},
			"b_flag": {
				Type:    "bool",
				NotNull: true,
			},
			"description": {
				Type:       "varchar",
				CharLength: 255,
			},
			"partition": {
				Type:         "int4",
				NumPrecision: 32,
			},
			"subpartition": {
				Type:         "float8",
				NumPrecision: 53,
			},
		},
		Indexes: map[string]schema.Index{
			"table_a_pkey": {
				Columns:  []string{"id"},
				IsUnique: true,
			},
			"name_idx": {
				Columns:  []string{"name"},
				IsUnique: true,
			},
			"partition_idx": {
				Columns: []string{"partition", "subpartition"},
			},
		},
	}

	validationResult, err := schema.Validate(ctx, db, expectedSchema, true)
	if err != nil {
		log.Fatal().Err(err).Msg("Validation failed")
	}
	ensureNoValidationErrors(validationResult)
}

func ensureNoValidationErrors(validationResult []string) {
	if len(validationResult) == 0 {
		return
	}
	for _, msg := range validationResult {
		log.Error().Msg(msg)
	}
	log.Fatal().Msg("DB Schema validation failed")
}
