package timeutil

import (
	"math/rand"
	"sync"
	"time"

	"github.com/leclerc04/go-tool/agl/base/concurrent"
	"github.com/leclerc04/go-tool/agl/util/errs"
)

// BackOff returns a exponential backoff duration between baseDelay and maxDelay.
func BackOff(retries int, baseDelay, maxDelay time.Duration) time.Duration {
	const multiplier = 1.3
	const randRatio = 0.4

	backOff, maxDelayF := float64(baseDelay), float64(maxDelay)
	for backOff < maxDelayF && retries > 0 {
		retries--
		backOff *= multiplier
	}
	if backOff > maxDelayF {
		backOff = maxDelayF
	}

	backOff -= rand.Float64() * randRatio * backOff
	if backOff < 0 {
		backOff = 0
	}
	return time.Duration(backOff)
}

// RetryUntil retries the action with computed backoff until timeout or action returns true.
// When action returns false, it should return error.
func RetryUntil(minDelay, maxDelay, timeOut time.Duration, action func() (ok bool, err error)) error {
	start := Now()

	var err error
	attempt := 0
	for Now().Sub(start) < timeOut {
		var ok bool
		ok, err = action()
		if ok {
			return err
		}
		delay := BackOff(attempt, minDelay, maxDelay)
		attempt++
		time.Sleep(delay)
	}
	if err == nil {
		// If you hit this path, your action is not correctly returning error.
		return errs.Newf("timeout retrying but there is no error specified.")
	}
	return err
}

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	concurrent.GoLite(func() {
		defer close(c)
		wg.Wait()
	})
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
