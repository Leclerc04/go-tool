package contextutil

import "context"

// Bool represents a context key for a boolean value.
type Bool Iface

// NewBool creates a new context key.
func NewBool() Bool {
	return Bool(NewIface())
}

// WithValue sets the value to context.
func (k Bool) WithValue(ctx context.Context, v bool) context.Context {
	return Iface(k).WithInterface(ctx, v)
}

// Value returns the value from context.
func (k Bool) Value(ctx context.Context) bool {
	v := Iface(k).Interface(ctx)
	if v == nil {
		return false
	}
	return v.(bool)
}
