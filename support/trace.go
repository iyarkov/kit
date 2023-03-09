package support

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func TraceNoOps() {
	provider := trace.NewNoopTracerProvider()
	Tracer = provider.Tracer("application")
}

func TraceConsole() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize stdout trace exporter")
	}

	provider := sdk.NewTracerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(Resource()),
	)
	otel.SetTracerProvider(provider)

	Tracer = provider.Tracer("tracer")

	log.Info().Msg("stdout trace exporter initialized")
	return func() {
		log.Info().Msg("Stopping stdout trace exporter")
		ctx := context.Background()

		if err := provider.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("stdout trace exporter shutdown failed")
		}
		log.Info().Msg("stdout trace exporter stopped")
	}
}
