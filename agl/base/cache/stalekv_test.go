package cache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/leclecr04/go-tool/agl/base/cache"
	"github.com/leclecr04/go-tool/agl/util/timeutil"
	"github.com/stretchr/testify/assert"
)

type counter struct {
	mu        sync.Mutex
	numCalled int
}

func (c *counter) Inc() {
	c.mu.Lock()
	c.numCalled++
	c.mu.Unlock()
}

func (c *counter) Value() int {
	c.mu.Lock()
	v := c.numCalled
	c.mu.Unlock()
	return v
}

type cacheKey struct {
	s string
}

func TestStaleKVEcho(t *testing.T) {
	ctx := context.Background()

	timeutil.SetAFakeTime()
	defer timeutil.UnSetFakeTime()

	c := counter{}
	injectError := false
	cache := cache.NewStaleKV("test", time.Minute, time.Minute*3, func(ctx context.Context, key interface{}) (interface{}, error) {
		c.Inc()
		if injectError {
			return "", fmt.Errorf("error %d", c.Value())
		}
		return "echo " + key.(cacheKey).s, nil
	})

	v, err := cache.Get(ctx, cacheKey{"foo"})
	assert.NoError(t, err)
	assert.Equal(t, "echo foo", v.(string))
	assert.Equal(t, 1, c.Value())

	v, err = cache.Get(ctx, cacheKey{"bar"})
	assert.NoError(t, err)
	assert.Equal(t, "echo bar", v.(string))
	assert.Equal(t, 2, c.Value())

	timeutil.AdvanceFakeTime(time.Minute + 1)
	for c.Value() == 2 {
		time.Sleep(time.Millisecond)
		_, err := cache.Get(ctx, cacheKey{"foo"})
		assert.NoError(t, err)
	}
	v, err = cache.Get(ctx, cacheKey{"foo"})
	assert.NoError(t, err)
	assert.Equal(t, "echo foo", v.(string))
	assert.Equal(t, 3, c.Value())

	timeutil.AdvanceFakeTime(time.Minute * 10)
	v, err = cache.Get(ctx, cacheKey{"baz"})
	assert.NoError(t, err)
	assert.Equal(t, "echo baz", v.(string))
	time.Sleep(time.Millisecond)
	assert.Equal(t, 1, cache.Size())

	injectError = true
	timeutil.AdvanceFakeTime(time.Minute + 1)
	for i := 0; i < 100; i++ {
		_, err = cache.Get(ctx, cacheKey{"foo"})
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond)
	}
	assert.Equal(t, "error 6", err.Error())
	// get again
	_, err = cache.Get(ctx, cacheKey{"foo"})
	assert.Equal(t, "error 7", err.Error())
	injectError = false
	v, err = cache.Get(ctx, cacheKey{"foo"})
	assert.NoError(t, err)
	assert.Equal(t, "echo foo", v.(string))
}
