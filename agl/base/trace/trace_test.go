package trace_test

import (
	"context"
	"testing"

	"github.com/leclecr04/go-tool/agl/base/trace"
)

func TestLog(t *testing.T) {
	t.Run("trace1", func(t *testing.T) {
		ctx := context.Background()
		ctx, done := trace.WithTrace(ctx, "Test")
		defer done()

		trace.Printf(ctx, "hello %s", "hi")
		trace.Printf(ctx, "hello %s", "hi2")
	})

	t.Run("trace2", func(t *testing.T) {
		ctx := context.Background()
		ctx, done := trace.WithTrace(ctx, "Test")
		defer done()

		trace.Printf(ctx, "hello %s", "hi")
		trace.Printf(ctx, "hello %s", "hi2")
	})
}
