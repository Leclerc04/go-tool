package timeutil

import (
	"context"
	"time"
)

// Sleep that supports context cancellation. Returns an cancelled error
// if it is canceled. Otherwise, no error.
func Sleep(ctx context.Context, dur time.Duration) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-time.After(dur):
		return nil
	}
}
