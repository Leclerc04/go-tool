package concurrent

import (
	"context"

	"github.com/leclerc04/go-tool/agl/base/trace"
	"github.com/leclerc04/go-tool/agl/util/errs"
	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	name string
	s    *semaphore.Weighted
}

func NewSemaphore(name string, n int64) *Semaphore {
	return &Semaphore{
		name: name,
		s:    semaphore.NewWeighted(n),
	}
}

func (sem *Semaphore) Acquire(ctx context.Context, n int64) error {
	if sem.s.TryAcquire(n) {
		return nil
	}
	trace.Printf(ctx, "Wait for semaphore: %s", sem.name)
	err := sem.s.Acquire(ctx, n)
	if err != nil {
		trace.Printf(ctx, "canceled")
		// Context cancel
		return errs.Wrap(err)
	}
	trace.Printf(ctx, "semaphore acquired")
	return nil
}

func (sem *Semaphore) Release(n int64) {
	sem.s.Release(n)
}
