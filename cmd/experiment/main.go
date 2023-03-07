package main

import (
	"context"
	"errors"
	"fmt"
	"microserver.rockyrunstream.com/foundation/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("Experiment")
	logger.Pretty()
	//logger.Json("host-1", "test app", "dev")
	//zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	{
		log.Trace().Msg("Trace")
		log.Debug().Str("AAA", "BBB").Msg("Debug")
		log.Log().Msg("Log")
		log.Info().Msg("Info")
		log.Warn().Msg("Warn")
		log.Error().Msg("Error")
		//log.Fatal().Msg("Fatal")
		//log.Panic().Msg("Panic")

	}

	{
		lg := log.With().Str("nnn", "mmmm").Logger()
		lg.Trace().Str("AAA", "BBB").Msg("Trace")
	}

	{
		ctx := context.WithValue(context.Background(), "foo", "bar")
		lg := log.With().Str("rq_id", "123").Logger()
		newContext := lg.WithContext(ctx)

		log.Ctx(newContext).Debug().Msg("MSSSS")
	}

	doWithStack()

	{
		err := fmt.Errorf("some err [%w] %s", fmt.Errorf("nested err [%w]: %s", fmt.Errorf("deep err"), "middle level"), "top level")
		log.Error().Msgf("Request failed: %v", err)
	}

	{
		log.Debug().Msg("With caller")
	}
	fmt.Println("Done")
}

func doWithStack() {
	log.Error().Stack().Err(errors.New("file open failed!")).Msg("Stack")
}

//func getLogger(ctx context.Context) *zerolog.Logger {
//	logFromCtx := logger.Ctx(ctx)
//	logFromCtx.With().Str("srv", "bb").Logger()
//	return logFromCtx
//}
