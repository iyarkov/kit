package telemetry

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var tracerProvider *sdk.TracerProvider

func WithSpan(ctx context.Context, name string, f func(context.Context) error) error {
	ctx, span := tracer.Start(ctx, name)
	defer span.End()
	return f(ctx)
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func traceNoOps() {
	provider := trace.NewNoopTracerProvider()
	tracer = provider.Tracer("application")
}

func traceConsole() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize stdout trace exporter")
	}

	tracerProvider = sdk.NewTracerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	tracer = tracerProvider.Tracer("tracer")

	log.Info().Msg("stdout trace exporter initialized")
}

func shutdownTrace() {
	if tracerProvider == nil {
		return
	}
	log.Info().Msg("Stopping trace exporter")
	ctx := context.Background()

	if err := tracerProvider.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("trace exporter shutdown failed")
	}
	log.Info().Msg("trace exporter stopped")

}

func traceDocker(ctx context.Context) {
	client := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure())
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize trace exporter")
	}

	tracerProvider = sdk.NewTracerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	tracer = tracerProvider.Tracer("tracer")

	log.Info().Msg("docker trace exporter initialized")
}

func traceFlush(ctx context.Context) {
	if tracerProvider == nil {
		return
	}
	err := tracerProvider.ForceFlush(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tracer flush failed")
	}
}
