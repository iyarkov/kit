package support

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSigTerm(hook func(os.Signal)) {
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sg := <-exit
		hook(sg)
	}()

}
