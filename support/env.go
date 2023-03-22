package support

import (
	"github.com/rs/zerolog/log"
)

type Manifest struct {
	Instance  string
	Name      string
	Version   string
	Namespace string
}

var AppManifest = Manifest{
	Instance:  "Node-1",
	Name:      "Experiment",
	Version:   "1.0.0.123",
	Namespace: "Dev",
}

func DumpAppInfo() {
	log.Info().Msg("Application info")
	log.Info().Str("Instance", AppManifest.Instance).Send()
	log.Info().Str("Name", AppManifest.Name).Send()
	log.Info().Str("Version", AppManifest.Version).Send()
	log.Info().Str("Namespace", AppManifest.Namespace).Send()
}
