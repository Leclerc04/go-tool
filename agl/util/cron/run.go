package cron

import (
	"context"

	"github.com/leclecr04/go-tool/agl/base/sentry"
	"github.com/leclecr04/go-tool/agl/base/trace"
	"github.com/leclecr04/go-tool/agl/util/must"
)

func runWithTraceAndSentry(ctx context.Context, familyName, traceName string, panicToSentry bool, action func(ctx context.Context) error) {
	ctx, done := trace.WithTrace(ctx, familyName+"/"+traceName)
	defer done()
	if panicToSentry {
		defer sentry.Recover(ctx, false)
	}
	must.Must(action(ctx))
}
