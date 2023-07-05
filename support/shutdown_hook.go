package support

import (
	"context"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"syscall"
)

var shutdownContext = WithContextId(context.Background(), uuid.NewString())

func OnSigTerm(hook func(context.Context, os.Signal)) {
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sg := <-exit
		hook(shutdownContext, sg)
	}()

}
