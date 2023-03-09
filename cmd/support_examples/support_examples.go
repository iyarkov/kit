package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"microserver.rockyrunstream.com/foundation/support"
	"time"
)

func main() {
	support.LoggerConsole()
	//support.LoggerCloud("host-1", "test app", "dev")

	support.MetricNoOps()
	//defer support.MetricConsole()()

	support.TraceNoOps()
	//defer support.TraceConsole()()

	support.DumpAppInfo()

	logs()
	metrics()
	traces()
}

func logs() {
	log.Trace().Msg("Trace")
	log.Debug().Str("AAA", "BBB").Msg("Debug")
	log.Log().Msg("Log")
	log.Info().Msg("Info")
	log.Warn().Msg("Warn")
	log.Error().Msg("Error")
}

func metrics() {
	log.Debug().Msg("Testing metrics")

	// Prepare
	counter, err := support.Meter.Int64Counter("some_counter")
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create some_counter")
	}
	attrs := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}

	// Use counter
	counter.Add(context.Background(), 1, attrs...)

	log.Debug().Str("z", "zz").Msg("Done testing")
}

func traces() {
	log.Debug().Msg("Testing traces")
	ctx := context.Background()
	ctx, topSpan := support.Tracer.Start(ctx, "top")
	for i := 0; i < 3; i++ {
		_, nestedSpan := support.Tracer.Start(ctx, fmt.Sprintf("nestedSpan %d", i))
		time.Sleep(time.Millisecond * 10)
		nestedSpan.End()

	}
	topSpan.End()
	log.Debug().Msg("Testing Done")
}
