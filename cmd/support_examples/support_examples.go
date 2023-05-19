//go:generate ./build_info.sh
package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/iyarkov/foundation/config"
	"github.com/iyarkov/foundation/logger"
	"github.com/iyarkov/foundation/schema"
	"github.com/iyarkov/foundation/support"
	"github.com/iyarkov/foundation/telemetry"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"os"
	"time"
)

type Configuration struct {
	App       support.Manifest
	Log       logger.Configuration
	Telemetry telemetry.Configuration
	Db        config.DbConfig
}

func main() {

	// ctx := context.Background()
	cfg := Configuration{}
	if err := config.Read(&cfg); err != nil {
		fmt.Printf("failed to read initial configuration %v", err)
		os.Exit(1)
	}
	cfg.App.Name = "Example"
	cfg.App.Version = version
	// Init App Info
	support.AppManifest = cfg.App

	// Init Logger
	logger.InitLogger(cfg.Log)
	telemetry.InitTelemetry(cfg.Telemetry)

	log.Info().Any("configuration", cfg).Send()

	// Connect to the DB
	logs()
	metrics()
	traces()
	db(cfg.Db)

	telemetry.Shutdown()
}

func logs() {
	log.Trace().Msg("Trace")
	log.Debug().Str("AAA", "BBB").Msg("Debug")
	log.Log().Msg("Log")
	log.Info().Msg("Info")
	log.Warn().Msg("Warn")
	log.Error().Msg("Error")
}

func metrics() {
	log.Debug().Msg("Testing metrics")

	// Prepare
	counter, err := telemetry.Meter.Int64Counter("some_counter")
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create some_counter")
	}
	attrs := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}

	// Use counter
	counter.Add(context.Background(), 1, attrs...)

	log.Debug().Str("z", "zz").Msg("Done testing")
}

func traces() {
	log.Debug().Msg("Testing traces")
	ctx := context.Background()
	ctx, topSpan := telemetry.Tracer.Start(ctx, "top")
	for i := 0; i < 3; i++ {
		_, nestedSpan := telemetry.Tracer.Start(ctx, fmt.Sprintf("nestedSpan %d", i))
		time.Sleep(time.Millisecond * 10)
		nestedSpan.End()

	}
	topSpan.End()
	log.Debug().Msg("Testing Done")
}

func db(cfg config.DbConfig) {
	ctx := context.Background()
	lg := zerolog.Ctx(ctx).With().Str("thread", "init").Logger()
	lg.Info().Msg("Initializing the database")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password.Value(),
		cfg.DbName)

	db, err := sql.Open("pgx", psqlconn)
	if err != nil {
		lg.Fatal().Err(err).Msg("DB connection failed")
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			lg.Warn().Err(closeErr).Msg("Failed to close the DB")
		}
	}()

	_, _, err = schema.Update(ctx, db, changeset)
	if err != nil {
		lg.Fatal().Err(err).Msg("DB schema upgrade failed")
	}
	messages, err := schema.Validate(ctx, db, expectedSchema, false)
	if err != nil {
		lg.Fatal().Err(err).Msg("DB schema validation failed")
	}
	if len(messages) != 0 {
		lg.Error().Msg("-------- DB schema invalid -------------")
		for _, m := range messages {
			lg.Error().Msgf("\t%s", m)
		}
		lg.Error().Msg("----------------------------------------")
		lg.Fatal().Msg("Application aborted")
	}
}
