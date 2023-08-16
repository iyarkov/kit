package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/iyarkov/kit/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"net/http"
	"os"
)

var Meter metric.Meter
var metricProvider *sdk.MeterProvider

func metricNoOps() {
	Meter = noop.NewMeterProvider().Meter("application")
}

func metricConsole(ctx context.Context) {
	log := zerolog.Ctx(ctx)
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
	Meter = metricProvider.Meter("meter")
	log.Info().Msg("stdout metric exporter initialized")
}

func shutdownMetric(ctx context.Context) {
	log := zerolog.Ctx(ctx)
	if metricProvider == nil {
		return
	}
	log.Info().Msg("Stopping metric exporter")

	if err := metricProvider.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("metric exporter shutdown failed")
	}
	log.Info().Msg("metric exporter stopped")
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

func metricDocker(ctx context.Context) {
	log := zerolog.Ctx(ctx)
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize prometheus metric exporter")
	}
	metricProvider = sdk.NewMeterProvider(sdk.WithReader(exporter))
	Meter = metricProvider.Meter(support.AppManifest.Name)

	go func() {
		srv := http.Server{
			Addr:    ":8081",
			Handler: promhttp.Handler(),
		}
		support.OnSigTerm(func(shutdownContext context.Context, signal os.Signal) {
			shutdownContext = logger.WithLogger(shutdownContext)
			shutdownLog := zerolog.Ctx(shutdownContext)
			shutdownLog.Info().Msg("Shutting down metrics HTTP server")
			if shutdownErr := srv.Shutdown(shutdownContext); shutdownErr != nil {
				shutdownLog.Error().Err(shutdownErr).Msg("Shutting down metrics HTTP server failed")
			} else {
				shutdownLog.Info().Msgf("metrics HTTP server stopped")
			}
		})
		listenErr := srv.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			zerolog.Ctx(ctx).Error().Err(listenErr).Msg("metrics HTTP server failed to start")
		}
	}()
}
