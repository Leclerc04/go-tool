package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leclecr04/go-tool/agl/base/mon"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/leclecr04/go-tool/agl/base/concurrent"
	"github.com/leclecr04/go-tool/agl/util/timeutil"
)

var (
	metricStaleKVSize = mon.NewGaugeVec(
		"cache", "stale_kv_size", "number of items in StaleKV", []string{"name"})
)

// StaleKV provides a local cache that could fetch value based on a given key.
// When a cached value became stale, a new get will trigger a new fetch, before the
// fetch complete, stale value is returned.
// A background GC will try to sweep the expired items. However, that is done best effort.
// So an expired value could still be returned.
type StaleKV struct {
	name   string
	loader func(ctx context.Context, key interface{}) (interface{}, error)
	stale  time.Duration
	expire time.Duration
	gcLife time.Duration

	mu      sync.RWMutex
	entries map[interface{}]*Local

	gcMu      sync.Mutex
	lastGC    time.Time
	gcRunning bool

	sizeGauge prometheus.Gauge
}

// NewSharedStaleKV creates a shared stale kv.
func NewSharedStaleKV(name string, stale time.Duration, expire time.Duration, loader func(ctx context.Context, key interface{}) (interface{}, error)) (ret *StaleKV) {
	func() {
		holder.mu.RLock()
		defer holder.mu.RUnlock()
		ret = holder.kv[name]
	}()
	if ret != nil {
		return ret
	}
	holder.mu.Lock()
	defer holder.mu.Unlock()
	ret = holder.kv[name]
	if ret != nil {
		return ret
	}
	ret = NewStaleKV(name, stale, expire, loader)
	holder.kv[name] = ret
	return ret
}

// NewStaleKV creates a StaleKV.
func NewStaleKV(name string, stale time.Duration, expire time.Duration, loader func(ctx context.Context, key interface{}) (interface{}, error)) *StaleKV {
	gcLife := expire
	if gcLife == 0 {
		gcLife = stale * 2
	}
	if gcLife < time.Minute {
		gcLife = time.Minute
	}
	if stale > expire {
		panic(fmt.Errorf("stale must be smaller than expire %v, %v", stale, expire))
	}
	return &StaleKV{
		name:    name,
		loader:  loader,
		stale:   stale,
		expire:  expire,
		gcLife:  gcLife,
		entries: make(map[interface{}]*Local),
		// Don't GC for the first minute.
		lastGC: timeutil.Now().Add(time.Minute),
		sizeGauge: metricStaleKVSize.With(mon.Labels{
			"name": name,
		}),
	}
}

// Get returns a value from the cache.
func (l *StaleKV) Get(ctx context.Context, key interface{}) (result interface{}, err error) {
	defer l.garbageCollect()
	l.mu.RLock()
	entry, ok := l.entries[key]
	l.mu.RUnlock()

	if !ok {
		entry = func() *Local {
			l.mu.Lock()
			defer l.mu.Unlock()
			entry = NewLocal(l.name, l.stale, func(ctx context.Context) (interface{}, error) {
				return l.loader(ctx, key)
			})
			l.entries[key] = entry
			l.sizeGauge.Set(float64(len(l.entries)))
			return entry
		}()
		return entry.Get(ctx)
	}

	if l.expire > 0 {
		loadedTime := entry.LoadedTime()
		expireTime := loadedTime.Add(l.expire)
		// if expire
		if !loadedTime.IsZero() && expireTime.After(timeutil.Now()) {
			return entry.ReloadIfExpire(ctx, expireTime)
		}
	}
	return entry.Get(ctx)
}

// garbageCollect performs garbage collection to remove items stale for twice of the time.
func (l *StaleKV) garbageCollect() {
	run := func() bool {
		l.gcMu.Lock()
		defer l.gcMu.Unlock()
		if l.gcRunning || l.lastGC.After(timeutil.Now().Add(-l.gcLife)) {
			return false
		}
		l.lastGC = timeutil.Now()
		l.gcRunning = true
		return true
	}()
	if !run {
		return
	}
	concurrent.GoLite(func() {
		var expires []interface{}
		cutoff := timeutil.Now().Add(-l.gcLife)
		func() {
			l.mu.RLock()
			defer l.mu.RUnlock()
			for k, entry := range l.entries {
				t := entry.LoadedTime()
				if t.IsZero() || t.After(cutoff) {
					continue
				}
				expires = append(expires, k)
			}
		}()
		if len(expires) > 0 {
			func() {
				l.mu.Lock()
				defer l.mu.Unlock()
				for _, k := range expires {
					delete(l.entries, k)
				}
				l.sizeGauge.Set(float64(len(l.entries)))
			}()
		}
		l.gcMu.Lock()
		l.gcRunning = false
		l.gcMu.Unlock()
	})
}

// Size returns number of values in the cache.
func (l *StaleKV) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}
