package contextutil

import "context"

// Int64 represents a context key for a int64 value.
type Int64 Iface

// NewInt64 creates a new context key.
func NewInt64() Int64 {
	return Int64(NewIface())
}

// WithValue sets the value to context.
func (k Int64) WithValue(ctx context.Context, v int64) context.Context {
	return Iface(k).WithInterface(ctx, v)
}

// Value returns the value from context.
func (k Int64) Value(ctx context.Context) int64 {
	v := Iface(k).Interface(ctx)
	if v == nil {
		return 0
	}
	return v.(int64)
}
