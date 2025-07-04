package concurrent

import "sync"

func NewStringSet() *StringSet {
	return &StringSet{
		sm: &sync.Map{},
	}
}

type StringSet struct {
	sm *sync.Map
}

func (c *StringSet) Contains(key string) bool {
	_, ok := c.sm.Load(key)
	return ok
}

func (c *StringSet) Delete(key string) {
	c.sm.Delete(key)
}

func (c *StringSet) Set(key string) {
	c.sm.Store(key, true)
}
