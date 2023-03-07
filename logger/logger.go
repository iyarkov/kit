package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func Pretty() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro}).
		With().
		Timestamp().
		Caller().
		Logger()
	zerolog.DefaultContextLogger = &log.Logger
}

func Json(host string, app string, env string) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Str("host", host).
		Str("app", app).
		Str("env", env).
		Logger()
	zerolog.DefaultContextLogger = &log.Logger

}
