package cache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/leclerc04/go-tool/agl/base/cache"
	"github.com/leclerc04/go-tool/agl/util/errs"
	"github.com/leclerc04/go-tool/agl/util/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestLocal(t *testing.T) {
	ctx := context.Background()

	timeutil.SetAFakeTime()
	defer timeutil.UnSetFakeTime()
	i := 0
	injectError := false
	l := cache.NewLocal("test", time.Minute, func(ctx context.Context) (interface{}, error) {
		i++
		if injectError {
			return 0, fmt.Errorf("error %d", i)
		}
		return i, nil
	})
	v, err := l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, v.(int))

	timeutil.AdvanceFakeTime(time.Minute + 10)
	// This still returns the old value, but it should trigger a load in background.
	v, err = l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, v.(int))

	// Spin until we get the new value.
	for v.(int) == 1 {
		v, err = l.Get(ctx)
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
	}
	v, err = l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, v.(int))

	// Force reload.
	l.Reload(ctx)
	v, err = l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, v.(int))

	injectError = true
	// make sure a new value is requested.
	timeutil.AdvanceFakeTime(time.Minute + 10)
	for i := 0; i < 100; i++ {
		_, err = l.Get(ctx)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond)
	}

	// it returns error.
	assert.Equal(t, "error 5", err.Error())

	// get again, since last one was having error.
	// the loader function will be reruned, regardless TTL.
	_, err = l.Get(ctx)
	assert.Equal(t, "error 6", err.Error())

	// When the error is clear, the value will be
	// retrieved immediately.
	injectError = false
	v, err = l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 7, v)

	v, err = l.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 7, v)
}

func TestLocalWithCancel(t *testing.T) {
	ch := make(chan string, 10)
	l := cache.NewLocal("test", time.Minute, func(ctx context.Context) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, errs.Wrap(ctx.Err())
		case v := <-ch:
			return v, nil
		}
	})

	{
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			time.Sleep(time.Millisecond * 10)
			cancel()
		}()

		v, err := l.Get(ctx)
		assert.True(t, errs.IsCancelled(err))
		assert.Nil(t, v)
	}

	ch <- "one"
	{
		ctx := context.Background()
		v, err := l.Get(ctx)
		assert.Nil(t, err)
		assert.Equal(t, v, "one")
	}

	{
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		v, err := l.Get(ctx)
		assert.Nil(t, err)
		assert.Equal(t, v, "one")
	}

	{
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		l.Reload(ctx)

		v, err := l.Get(ctx)
		assert.Nil(t, err)
		assert.Equal(t, v, "one")
	}
}
