package support

import (
	"context"
	"github.com/rs/zerolog"
	"reflect"
)

type Closeable interface {
	Close() error
}

func CloseWithWarning(ctx context.Context, cls Closeable, msg string) {
	if cls == nil || reflect.ValueOf(cls).IsNil() {
		return
	}
	if err := cls.Close(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg(msg)
	}
}

func DoWithWarning(ctx context.Context, callback func() error, msg string) {
	if callback == nil || reflect.ValueOf(callback).IsNil() {
		return
	}
	if err := callback(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg(msg)
	}
}
