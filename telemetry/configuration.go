package telemetry

import (
	"github.com/iyarkov/foundation/support"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type Configuration struct {
	Mode string
}

func InitTelemetry(config Configuration) {
	switch config.Mode {
	case "console":
		log.Debug().Msg("Initializing console telemetry")
		metricConsole()
		traceConsole()
	default:
		log.Debug().Msg("Initializing NoOps telemetry")
		metricNoOps()
		traceNoOps()
	}
}

func Shutdown() {
	shutdownMetric()
	shutdownConsole()
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
