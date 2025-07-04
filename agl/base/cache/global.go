package cache

import (
	"sync"
)

var holder cacheHolder

func init() {
	ForgetSharedLocalCaches()
}

type cacheHolder struct {
	mu     sync.RWMutex
	single map[string]*Local
	kv     map[string]*StaleKV
}

// ForgetSharedLocalCaches forget all shared local caches created.
func ForgetSharedLocalCaches() {
	holder.mu.Lock()
	defer holder.mu.Unlock()
	holder.single = make(map[string]*Local)
	holder.kv = make(map[string]*StaleKV)
}
