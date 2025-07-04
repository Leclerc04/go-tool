package ctxstore

import (
	"context"
	"sync"

	dataloader "github.com/leclecr04/go-tool/agl/util/dataloader"
)

// ctxstore manages a per context dataloader.
// This is useful to make loading data more efficiently (by concurrency and deduplication)
// for a given request.

type contextKey struct{}

type root struct {
	sch *dataloader.Scheduler

	mu  sync.Mutex
	dlm map[string]*dataloader.DataLoader
}

// StartRoot runs the function with context that can share data loading.
func StartRoot(ctx context.Context, f func(ctx context.Context)) {
	dataloader.RunWithScheduler(func(sch *dataloader.Scheduler) {
		ctx := context.WithValue(ctx, contextKey{}, &root{
			sch: sch,
			dlm: make(map[string]*dataloader.DataLoader),
		})
		f(ctx)
	})
}

// Scheduler returns the scheduler installed on the context or nil.
func Scheduler(ctx context.Context) *dataloader.Scheduler {
	rootPtr := ctx.Value(contextKey{})
	if rootPtr == nil {
		return nil
	}
	return rootPtr.(*root).sch
}

// EnsureLoader ensures a loader is defined for the current context.
func EnsureLoader(ctx context.Context, loaderName string, loaderFunc func(keys []interface{}) []dataloader.Value) *dataloader.DataLoader {
	rootPtr := ctx.Value(contextKey{})
	if rootPtr == nil {
		return dataloader.New(nil, loaderFunc)
	}
	root := rootPtr.(*root)
	root.mu.Lock()
	defer root.mu.Unlock()
	dl, ok := root.dlm[loaderName]
	if ok {
		return dl
	}
	dl = dataloader.New(root.sch, loaderFunc)
	root.dlm[loaderName] = dl
	return dl
}

// Load loads a data with the loader on the current context.
func Load(ctx context.Context, loaderName string, loaderFunc func(keys []interface{}) []dataloader.Value, key interface{}) dataloader.Value {
	return EnsureLoader(ctx, loaderName, loaderFunc).Load(key)
}

// LoadSimple loads a data with the loader on the current context. It just simply parallelize the load.
func LoadSimple(ctx context.Context, loaderName string, loaderFunc func(key interface{}) dataloader.Value, key interface{}) dataloader.Value {
	return Load(ctx, loaderName, dataloader.Parallel(loaderFunc), key)
}
