package support

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func LoggerConsole() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro}).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.DebugLevel)
	zerolog.DefaultContextLogger = &log.Logger
}

func LoggerCloud() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Str("instance", AppInfo.Instance).
		Str("version", AppInfo.Version).
		Str("app", AppInfo.Name).
		Str("namespace", AppInfo.Namespace).
		Logger().
		Level(zerolog.InfoLevel)
	zerolog.DefaultContextLogger = &log.Logger

}
