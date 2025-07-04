package concurrent_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/leclecr04/go-tool/agl/base/concurrent"
	"github.com/leclecr04/go-tool/agl/util/must"
)

func TestKeyedMutex(t *testing.T) {
	km := concurrent.NewKeyedMutex()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := context.Background()
		t.Log(time.Now(), "r1 Lock hello1")
		unlock, err := km.Lock(ctx, "hello1")
		defer unlock(ctx)
		must.Must(err)
		t.Log(time.Now(), "r1 Locked hello1")
		time.Sleep(time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()
		t.Log(time.Now(), "r2 Lock hello2")
		unlock, err := km.Lock(ctx, "hello2")
		defer unlock(ctx)
		must.Must(err)
		t.Log(time.Now(), "r2 Locked hello2")
		time.Sleep(time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()
		t.Log(time.Now(), "r3 Lock hello1")
		unlock, err := km.Lock(ctx, "hello1")
		defer unlock(ctx)
		must.Must(err)
		t.Log(time.Now(), "r3 Locked hello1")
		time.Sleep(time.Second)
	}()

	wg.Wait()
}
