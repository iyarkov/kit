package support

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type Manifest struct {
	Instance  string
	Name      string
	Version   string
	Namespace string
}

var AppInfo = Manifest{
	Instance:  "Node-1",
	Name:      "Experiment",
	Version:   "1.0.0.123",
	Namespace: "Dev",
}

func DumpAppInfo() {
	log.Info().Msg("Application info")
	log.Info().Str("Instance", AppInfo.Instance).Send()
	log.Info().Str("Name", AppInfo.Name).Send()
	log.Info().Str("Version", AppInfo.Version).Send()
	log.Info().Str("Namespace", AppInfo.Namespace).Send()
}

func Resource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceInstanceID(AppInfo.Instance),
		semconv.ServiceName(AppInfo.Name),
		semconv.ServiceVersion(AppInfo.Version),
		semconv.ServiceNamespace(AppInfo.Namespace),
	)
}
