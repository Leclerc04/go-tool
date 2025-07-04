package cache

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/leclecr04/go-tool/agl/base/concurrent"
	"github.com/leclecr04/go-tool/agl/base/sentry"
	"github.com/leclecr04/go-tool/agl/base/trace"
	"github.com/leclecr04/go-tool/agl/util/errs"
	"github.com/leclecr04/go-tool/agl/util/timeutil"
)

// Local provides a local in memory cache that loads asynchronously.
type Local struct {
	loader func(ctx context.Context) (interface{}, error)
	ttl    time.Duration
	name   string

	mu         sync.Mutex
	v          atomic.Value // holds a localV.
	refreshing atomic.Bool
}

type localV struct {
	v      interface{} // holds the result of loader.
	err    error
	panicV error
	loaded time.Time
}

// NewSharedLocal allocates a new Local cache, and store it in global process space.
// If repeatedly call with the same name, the previous one will be returned.
// The name should be uniquely identifying the content.
func NewSharedLocal(name string, ttl time.Duration, loader func(ctx context.Context) (interface{}, error)) (l *Local) {
	func() {
		holder.mu.RLock()
		defer holder.mu.RUnlock()
		l = holder.single[name]
	}()
	if l != nil {
		return l
	}
	holder.mu.Lock()
	defer holder.mu.Unlock()
	l = holder.single[name]
	if l != nil {
		return l
	}
	l = NewLocal(name, ttl, loader)
	holder.single[name] = l
	return l
}

// NewLocal creates a new instance of Local. The provided loader func is used
// to retrieve the value to be cached. If it does not return error, the
// result will be cached. The cache will be refresh asynchronously when
// the next Get is called after the TTL.
// If the loader func returns error, this error will be propagated to the Get
// result, and the error is not cached, so the next Get will try to fetch the
// value synchronously.
func NewLocal(name string, ttl time.Duration, loader func(ctx context.Context) (interface{}, error)) *Local {
	l := &Local{
		name:   name,
		loader: loader,
		ttl:    ttl,
		mu:     sync.Mutex{},
	}
	return l
}

// Get returns the cached value.
func (l *Local) Get(ctx context.Context) (interface{}, error) {
	return l.get(ctx)
}

func (l *Local) get(ctx context.Context) (interface{}, error) {
	v0 := l.v.Load()

	var v *localV
	if v0 == nil {
		var err error
		v, err = l.loadInternal(ctx, false, future)
		if err != nil {
			return nil, err
		}
	} else {
		v = v0.(*localV)
	}

	if v.err != nil || v.panicV != nil {
		// if there is recorded error
		var err error
		v, err = l.loadInternal(ctx, false, future)
		if err != nil {
			return nil, err
		}
		return unpackLocalV(v)
	}
	if v.loaded.Before(timeutil.Now().Add(-l.ttl)) {
		if l.refreshing.CompareAndSwap(false, true) {
			// if expire
			concurrent.Go(ctx, l.name, func(ctx context.Context) error {
				defer l.refreshing.Store(false)
				_, err := l.loadInternal(ctx, false, future)
				if err != nil {
					trace.Printf(ctx, "cache refresh error: ", err)
				}
				return nil
			})
		}
	}
	return unpackLocalV(v)
}

func unpackLocalV(v *localV) (interface{}, error) {
	if v.panicV != nil {
		panic(v.panicV)
	}
	if v.err != nil {
		return nil, v.err
	}
	return v.v, nil
}

// Reload forces cache reload, blocks until the new value is loaded.
func (l *Local) Reload(ctx context.Context) {
	_, err := l.loadInternal(ctx, true, future)
	if err != nil {
		trace.Printf(ctx, "cache refresh error: ", err)
	}
}

var future = time.Date(21002, 9, 16, 19, 17, 23, 0, time.UTC)

func (l *Local) ReloadIfExpire(ctx context.Context, expiration time.Time) (interface{}, error) {
	v, err := l.loadInternal(ctx, false, expiration)
	if err != nil {
		return nil, err
	}
	return unpackLocalV(v)
}

// LoadedTime returns when the current value is loaded, or zero if never.
func (l *Local) LoadedTime() time.Time {
	v := l.v.Load()
	if v == nil {
		return time.Time{}
	}
	return v.(*localV).loaded
}

func (l *Local) loadInternal(ctx context.Context, force bool, expiration time.Time) (*localV, error) {
	if ctx.Err() != nil {
		return nil, errs.Wrap(ctx.Err())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if !force {
		v := l.v.Load()
		// if there is value and there is no error, panic or expiration
		if v != nil {
			lv := v.(*localV)
			now := timeutil.Now()
			loaded := lv.loaded
			if lv.err == nil && lv.panicV == nil && loaded.After(now.Add(-l.ttl)) && loaded.Before(expiration) {
				return lv, nil
			}

		}
	}
	lv := &localV{}
	err := func() (errOut error) {
		// if panic, record error
		defer func() {
			panicV := recover()
			if panicV != nil {
				if panicErr, ok := panicV.(error); ok && errs.IsCancelled(panicErr) {
					errOut = panicErr
					return
				}
				bs := make([]byte, 4096)
				runtime.Stack(bs, false)
				err := errs.Newf("recover from panic: %v, stack: %s", panicV, string(bs))
				sentry.Error(ctx, err)
				lv.panicV = err
			}
		}()
		// if loader returns error，record error but keep value empty
		v, err := l.loader(ctx)
		if err != nil {
			if errs.IsCancelled(ctx.Err()) {
				// Cancellation is special, we don't remember it, we treat as such load never happened.
				return ctx.Err()
			}
			lv.err = err
			return nil
		}
		lv.v = v
		return nil
	}()
	if err != nil {
		return nil, err
	}
	// store lv
	lv.loaded = timeutil.Now()
	l.v.Store(lv)
	return lv, nil
}
