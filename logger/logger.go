package logger

import (
	"context"
	"github.com/iyarkov/foundation/support"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Mode int

type Configuration struct {
	Mode  string
	Level string
}

func InitLogger(config Configuration) {
	var level zerolog.Level
	switch config.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel
	default:
		level = zerolog.InfoLevel
	}

	if config.Mode == "cloud" {
		zerolog.TimeFieldFormat = time.RFC3339
		log.Logger = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Caller().
			Str("instance", support.AppManifest.Instance).
			Str("version", support.AppManifest.Version).
			Str("app", support.AppManifest.Name).
			Str("namespace", support.AppManifest.Namespace).
			Logger().
			Level(level)
		zerolog.DefaultContextLogger = &log.Logger
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro}).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger().
			Level(level)

		zerolog.DefaultContextLogger = &log.Logger
	}
	log.Info().Msg("logger system initialized")
}

func SetContextLogger(ctx context.Context) context.Context {
	idx := support.GetContextId(ctx)
	if idx != "" {
		return log.Logger.With().Str("contextId", idx).Logger().WithContext(ctx)
	}
	return ctx
}
