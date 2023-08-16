package telemetry

import (
	"context"
	"github.com/iyarkov/kit/logger"
	"github.com/iyarkov/kit/support"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"os"
)

type Configuration struct {
	Mode string
}

func InitTelemetry(ctx context.Context, cfg *Configuration) {
	log := zerolog.Ctx(ctx)
	log.Debug().Msgf("InitTelemetry: %s", cfg.Mode)
	switch cfg.Mode {
	case "console":
		log.Debug().Msg("Initializing console telemetry")
		metricConsole(ctx)
		traceConsole(ctx)
	case "docker":
		log.Debug().Msg("Initializing docker telemetry")
		metricDocker(ctx)
		traceDocker(ctx)
	default:
		log.Debug().Msg("Initializing NoOps telemetry")
		metricNoOps()
		traceNoOps()
	}
	support.OnSigTerm(func(shutdownCtx context.Context, signal os.Signal) {
		shutdownCtx = logger.WithLogger(shutdownCtx)
		shutdownMetric(shutdownCtx)
		shutdownTrace(shutdownCtx)
	})
}

func Flush(ctx context.Context) {
	metricFlush(ctx)
	traceFlush(ctx)
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceInstanceID(support.AppManifest.Instance),
		semconv.ServiceName(support.AppManifest.Name),
		semconv.ServiceVersion(support.AppManifest.Version),
		semconv.ServiceNamespace(support.AppManifest.Namespace),
	)
}
