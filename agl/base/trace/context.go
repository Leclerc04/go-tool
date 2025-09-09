package trace

import (
	"context"
	"strings"
	"time"

	"github.com/leclerc04/go-tool/agl/util/errs"

	"github.com/golang/glog"
)

type contextKey struct{}

// Region annotates the start and end of a some process.
func Region(ctx context.Context, msg string, kv ...interface{}) func(error) {
	maybe := ctx.Value(contextKey{})
	if maybe == nil {
		return func(error) {}
	}
	start := time.Now()
	Printc(ctx, msg, kv...)
	tr := maybe.(T)
	tr.Indent(1)
	done := false
	return func(err error) {
		if done {
			return
		}
		done = true
		Printc(ctx, "done", "err", err, "duration", time.Since(start))
		tr.Indent(-1)
	}
}

func Printc(ctx context.Context, msg string, kv ...interface{}) {
	kvStr := errs.ParamKVToString(errs.BuildParamKV(kv...))
	sep := ""
	if kvStr != "" && msg != "" {
		sep = "\t"
	}
	Printf(ctx, "%s%s%s", msg, sep, kvStr)
}

// Printf prints to the trace of current context.
func Printf(ctx context.Context, msg string, args ...interface{}) {
	maybe := ctx.Value(contextKey{})
	if maybe == nil {
		return
	}
	tr := maybe.(T)
	tr.Printf(msg, args...)
}

// VLogf prints to the trace of current context and logs to V(1) log.
func VLogf(ctx context.Context, msg string, args ...interface{}) {
	Printf(ctx, msg, args...)
	glog.V(1).Infof(msg, args...)
}

// WithTrace creates a new context by attaching a trace.
func WithTrace(ctx context.Context, name string) (context.Context, func()) {
	maybe := ctx.Value(contextKey{})
	if maybe != nil {
		tr := maybe.(T)
		tr.Printf("sub-trace: %s", name)
	}

	m := strings.SplitN(name, "/", 2)
	family := m[0]
	method := "run"
	if len(m) == 2 {
		method = m[1]
	}
	tr, done := New(family, method)
	return context.WithValue(ctx, contextKey{}, tr), done
}

// WithGLogTrace creates a new context attacing a glog trace.
func WithGLogTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, glogTrace)
}

// tr returns trace attached to the context.
func tr(ctx context.Context) T {
	maybe := ctx.Value(contextKey{})
	if maybe == nil {
		// TODO: create a special trace for this.
		return Noop
	}
	return maybe.(T)
}

// MarkFailed marks the current trace as failed.
func MarkFailed(ctx context.Context) {
	tr := tr(ctx)
	tr.SetError()
}

// Scoped is a helper function to record the start and end of a function.
// Usage:
//
//	func dbRead(ctx context.Context) (err error) {
//		defer Scoped(ctx, "db read")(&err)
//		return someDbReadOp()
//	}
//
// The above will print:
// 1. db read at the begining of the function.
// 2. done at the end.
// 3. error: <err>, if error is not nil.
func Scoped(ctx context.Context, description string) func(perr *error) {
	rdone := Region(ctx, description)
	return func(perr *error) {
		if perr == nil || *perr == nil {
			rdone(nil)
			return
		}
		rdone(*perr)
	}
}
