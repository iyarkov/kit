package telemetry

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer
var tracerProvider *sdk.TracerProvider

func traceNoOps() {
	provider := trace.NewNoopTracerProvider()
	Tracer = provider.Tracer("application")
}

func traceConsole() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize stdout trace exporter")
	}

	tracerProvider = sdk.NewTracerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	Tracer = tracerProvider.Tracer("tracer")

	log.Info().Msg("stdout trace exporter initialized")
	return func() {
	}
}

func shutdownConsole() {
	if tracerProvider == nil {
		return
	}
	log.Info().Msg("Stopping stdout trace exporter")
	ctx := context.Background()

	if err := tracerProvider.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("stdout trace exporter shutdown failed")
	}
	log.Info().Msg("stdout trace exporter stopped")

}
