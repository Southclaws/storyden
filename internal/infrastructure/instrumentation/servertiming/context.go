package servertiming

import "context"

type contextKey struct{}

var key contextKey = struct{}{}

func WithCollector(ctx context.Context, collector *Collector) context.Context {
	return context.WithValue(ctx, key, collector)
}

func FromContext(ctx context.Context) *Collector {
	if collector, ok := ctx.Value(key).(*Collector); ok {
		return collector
	}
	return nil
}
