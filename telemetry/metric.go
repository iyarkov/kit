package telemetry

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"os"
)

var Meter metric.Meter
var metricProvider *sdk.MeterProvider

func metricNoOps() {
	Meter = noop.NewMeterProvider().Meter("application")
}

func metricConsole() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	exporter, err := stdoutmetric.New(stdoutmetric.WithEncoder(encoder))
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize stdout metric exporter")
	}

	// Register the exporter with an SDK via a periodic reader.
	metricProvider = sdk.NewMeterProvider(
		sdk.WithResource(newResource()),
		sdk.WithReader(sdk.NewPeriodicReader(exporter)),
	)
	Meter = metricProvider.Meter(
		"meter",
		metric.WithInstrumentationVersion("1.0"),
		metric.WithSchemaURL("my.app.com"),
	)
	log.Info().Msg("stdout metric exporter initialized")
}

func shutdownMetric() {
	if metricProvider == nil {
		return
	}
	log.Info().Msg("Stopping stdout metric exporter")
	ctx := context.Background()

	if err := metricProvider.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("stdout metric exporter shutdown failed")
	}
	log.Info().Msg("stdout metric exporter stopped")
}

func metricFlush(ctx context.Context) {
	if metricProvider == nil {
		return
	}
	err := metricProvider.ForceFlush(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tracer flush failed")
	}
}
