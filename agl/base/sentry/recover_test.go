package sentry_test

import (
	"context"
	"testing"

	"github.com/leclecr04/go-tool/agl/base/sentry"
)

func TestRecover(t *testing.T) {
	func() {
		defer sentry.Recover(context.Background(), false)
		panic("hi")
	}()
}
