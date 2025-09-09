package sentry

import (
	"context"

	"github.com/leclerc04/go-tool/agl/util/errs"
)

// Recover recovers panic and report it to sentry. And re-panic.
// Must be called directly by a defer, like this: defer sentry.Recover(false)
func Recover(ctx context.Context, repanic bool) {
	if r := recover(); r != nil {
		ErrorDepth(ctx, 0, errs.Newf("Recover panic: %v", r))
		if repanic {
			panic(r)
		}
	}
}

// RecoverAndSetError recovers panic and report it to sentry. And set err.
func RecoverAndSetError(ctx context.Context, err *error) {
	if r := recover(); r != nil {
		*err = errs.Newf("Recover panic: %v", r)
		ErrorDepth(ctx, 1, *err)
	}
}
