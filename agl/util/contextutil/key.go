package contextutil

import (
	"context"
	"sync/atomic"
)

var nextkey uint32

// Bool represents a context key for untyped value.
type Iface struct {
	v uint32
}

// NewIface creates a new context key holding any untyped value.
func NewIface() Iface {
	atomic.AddUint32(&nextkey, 1)
	return Iface{
		v: nextkey,
	}
}

// WithInterface sets the value to context.
func (k Iface) WithInterface(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, k, v)
}

// Interface returns the value from context.
func (k Iface) Interface(ctx context.Context) interface{} {
	return ctx.Value(k)
}
