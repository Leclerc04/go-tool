package cmdline

import "context"

type MethodRegistry struct {
	methods map[string]func(ctx context.Context) error
}

func NewMethodRegistry() *MethodRegistry {
	return &MethodRegistry{
		methods: map[string]func(ctx context.Context) error{},
	}
}

func (r *MethodRegistry) Register(name string, f func(ctx context.Context) error) {
	r.methods[name] = f
}

func (r *MethodRegistry) Run(ctx context.Context, name string) error {
	f := r.methods[name]
	if f == nil {
		panic("Unknown method " + name)
	}
	return f(ctx)
}
