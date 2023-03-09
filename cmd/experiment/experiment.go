package main

import (
	"github.com/rs/zerolog/log"
	"microserver.rockyrunstream.com/foundation/support"
)

func main() {
	support.LoggerConsole()
	//support.LogsJson("host-1", "test app", "dev")
	support.MetricConsole()
	//support.MetricsCloud("host-1", "test app", "dev")
	{
		log.Trace().Msg("Trace")
		log.Debug().Str("AAA", "BBB").Msg("Debug")
		log.Log().Msg("Log")
		log.Info().Msg("Info")
		log.Warn().Msg("Warn")
		log.Error().Msg("Error")
	}
	log.Debug().Msg("Testing timer")
}
