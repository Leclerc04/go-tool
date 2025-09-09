package cron

import (
	"context"
	"math/rand"
	"time"

	"github.com/leclerc04/go-tool/agl/base/concurrent"
)

// RunPeriodically runs the given action periodically in a new goroutine.
// Returns a function used to stop.
func RunPeriodically(ctx context.Context, runOnce bool, interval time.Duration, name string, action func(ctx context.Context) error) func() {
	if runOnce {
		runWithTraceAndSentry(ctx, "runonce."+name, "run", runOnce, action)
		return func() {
		}
	}
	quit := make(chan struct{})
	concurrent.Go(ctx, "periodic_loop."+name, func(ctx context.Context) error {
		if interval > time.Second {
			// Randomly delay the first fire. So different RunPeriodically with the same interval would start at different time.
			select {
			case <-time.After(time.Duration(interval.Seconds() * rand.Float64() * float64(time.Second))):
				runWithTraceAndSentry(ctx, "periodic."+name, "run", true, action)
			case <-quit:
				return nil
			}
		}
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				runWithTraceAndSentry(ctx, "periodic."+name, "run", true, action)
			case <-quit:
				ticker.Stop()
				return nil
			}
		}
	})

	return func() {
		quit <- struct{}{}
	}
}
