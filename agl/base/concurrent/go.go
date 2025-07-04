package concurrent

import (
	"context"
	"sync"

	"github.com/leclecr04/go-tool/agl/base/mon"
	"github.com/leclecr04/go-tool/agl/base/sentry"
	"github.com/leclecr04/go-tool/agl/base/trace"
)

var (
	metricInflightGo = mon.NewGauge("concurrent", "inflight_go", "number of goroutines active started by concurrent.Go")
)

// GoRoutine holds the information of a running go routine.
type GoRoutine struct {
	cancelCtx func()

	cancelOnce sync.Once
	finalErr   error
	finished   sync.WaitGroup
}

// Cancel cancel the context of the go routine.
func (gr *GoRoutine) Cancel() {
	gr.cancelOnce.Do(gr.cancelCtx)
}

// Wait waits for the goroutine to be finished.
func (gr *GoRoutine) Wait(ctx context.Context) error {
	trace.Printf(ctx, "waiting for go routine")
	gr.finished.Wait()
	trace.Printf(ctx, "waiting for go routine done")
	return gr.finalErr
}

// Go starts a new goroutine. If it panics, catch the error and send it to sentry.
// It returns a object for control. It is safe to ignore it.
func Go(ctx context.Context, name string, f func(ctx context.Context) error) *GoRoutine {
	asyncCtx := sentry.CreateDetachedContext(ctx)
	asyncCtx, cancelAsyncCtx := context.WithCancel(asyncCtx)
	gr := &GoRoutine{
		cancelCtx: cancelAsyncCtx,
	}
	gr.finished.Add(1)
	trace.Printc(ctx, "started go routine", "name", name)
	go func() {
		defer gr.finished.Done()
		defer gr.Cancel()
		defer mon.TrackInflight(metricInflightGo)()
		ctx, done := trace.WithTrace(asyncCtx, name)
		defer done()
		defer sentry.RecoverAndSetError(ctx, &gr.finalErr)

		err := f(ctx)
		if err != nil {
			sentry.Error(ctx, err)
			gr.finalErr = err
		}
	}()
	return gr
}

// GoLite is a lite version of Go. Prefer not using this when in doubt.
func GoLite(f func()) {
	go func() {
		defer mon.TrackInflight(metricInflightGo)()
		defer sentry.Recover(context.TODO(), false)
		f()
	}()
}

// GoChild starts a goroutine, and it requires the current goroutine to wait for its finish.
func GoChild(ctx context.Context, f func(ctx context.Context) error) func() error {
	errChan := make(chan error)
	go func() {
		errChan <- f(ctx)
	}()
	return func() error {
		return <-errChan
	}
}
