package telemetry

import (
	"context"
	"github.com/iyarkov/foundation/support"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
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
		metricConsole()
		traceConsole()
	case "docker":
		log.Debug().Msg("Initializing docker telemetry")
		metricNoOps()
		traceDocker(ctx)
	default:
		log.Debug().Msg("Initializing NoOps telemetry")
		metricNoOps()
		traceNoOps()
	}
}

func Shutdown() {
	shutdownMetric()
	shutdownTrace()
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
