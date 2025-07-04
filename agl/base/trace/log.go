package trace

import (
	"context"
	"fmt"

	"github.com/golang/glog"
)

// Debugf vlog and print to the trace.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	Printf(ctx, format, a...)
	if glog.V(1) {
		glog.InfoDepth(2, fmt.Sprintf(format, a...))
	}
}

// Infof log and print to the trace.
func Infof(ctx context.Context, format string, a ...interface{}) {
	Printf(ctx, format, a...)
	glog.InfoDepth(2, fmt.Sprintf(format, a...))
}

// Errorf log and print to the trace.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	Printf(ctx, format, a...)
	glog.ErrorDepth(2, fmt.Sprintf(format, a...))
}

// Warningf log and print to the trace.
func Warningf(ctx context.Context, format string, a ...interface{}) {
	Printf(ctx, format, a...)
	glog.WarningDepth(2, fmt.Sprintf(format, a...))
}
