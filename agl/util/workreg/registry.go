package onceutil

import (
	"context"
	"sync"

	"github.com/leclerc04/go-tool/agl/base/concurrent"
	"github.com/leclerc04/go-tool/agl/base/sentry"
)

// workreg provides a way to run things once.

// Registry allows some work to be registered and run.
type Registry struct {
	mu        sync.Mutex
	nextID    int
	onceWorks map[int]*OnceWork
}

// NewRegistry creates a registry.
func NewRegistry() *Registry {
	return &Registry{
		onceWorks: map[int]*OnceWork{},
	}
}

// OnceWork is a piece of work that other logic depends on.
type OnceWork struct {
	name   string
	f      func(ctx context.Context) error
	result error
	once   sync.Once
}

// Run runs the work. if it is already run, return the stored error.
func (w *OnceWork) Run(ctx context.Context) error {
	w.once.Do(func() {
		var err error
		defer sentry.RecoverAndSetError(ctx, &err)
		err = w.f(ctx)
		w.result = err
	})
	return w.result
}

// NewOnceWork defines a work that only run once.
func (r *Registry) NewOnceWork(name string, f func(ctx context.Context) error) *OnceWork {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	w := &OnceWork{
		name: name,
		f:    f,
	}
	r.onceWorks[id] = w
	r.nextID++
	return w
}

// Run runs given works in parallel.
func Run(ctx context.Context, works ...*OnceWork) error {
	var dones []func() error
	for _, w := range works {
		dones = append(dones, concurrent.GoChild(ctx, w.f))
	}
	var err error
	for _, done := range dones {
		dErr := done()
		if err == nil {
			err = dErr
		}
	}
	return err
}
