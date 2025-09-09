package concurrent

import (
	"context"
	"sync"

	"github.com/leclerc04/go-tool/agl/util/errs"
	semaphore "golang.org/x/sync/semaphore"
)

type KeyedMutex struct {
	mus map[interface{}]*kmItem
	mu  sync.Mutex
}

type kmItem struct {
	sem *semaphore.Weighted
	ref int
}

func NewKeyedMutex() *KeyedMutex {
	return &KeyedMutex{
		mus: make(map[interface{}]*kmItem),
	}
}

func (km *KeyedMutex) unref(key interface{}, kmI *kmItem) {
	km.mu.Lock()
	defer km.mu.Unlock()

	kmI.ref -= 1
	if kmI.ref > 0 {
		return
	}
	delete(km.mus, key)
}

// Locks a "mutex" of given key, respects context cancellation.
// The returned callback must be called to release the lock.
func (km *KeyedMutex) Lock(ctx context.Context, key interface{}) (func(ctx context.Context), error) {
	km.mu.Lock()
	m, ok := km.mus[key]
	if !ok {
		m = &kmItem{}
		m.sem = semaphore.NewWeighted(1)
		km.mus[key] = m
	}
	m.ref += 1
	km.mu.Unlock()

	err := m.sem.Acquire(ctx, 1)
	if err != nil {
		km.unref(key, m)
		// Context is cancelled.
		return func(context.Context) {
		}, errs.Wrap(err)
	}

	return func(ctx context.Context) {
		m.sem.Release(1)
		km.unref(key, m)
	}, nil
}
