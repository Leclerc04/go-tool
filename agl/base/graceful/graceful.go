package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/golang/glog"
)

// CancelWhenTerminated returns a context that will be cancel when a termination signal is received.
func CancelWhenTerminated(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	var cancelOnce sync.Once
	signals := make(chan os.Signal, 1)

	runCancel := func() {
		cancelOnce.Do(func() {
			signal.Stop(signals)
			cancel()
		})
	}

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		glog.Infof("%v signal received canceling context.", sig)
		runCancel()
		for range signals {
		}
	}()
	return ctx, func() {
		runCancel()
	}
}
