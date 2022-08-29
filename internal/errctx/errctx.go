package errctx

import (
	"context"
	"errors"
)

type contextKey struct{}

type withContext struct {
	underlying error
	meta       map[string]string
}

func (e *withContext) Error() string  { return e.underlying.Error() }
func (e *withContext) Cause() error   { return e.underlying }
func (e *withContext) Unwrap() error  { return e.underlying }
func (e *withContext) String() string { return e.Error() }

// WithMeta wraps a context with some arbitrary string based key-value metadata.
func WithMeta(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 != 0 {
		panic("odd number of key-value pair arguments")
	}

	var data map[string]string

	// overwrite any existing context metadata
	if meta, ok := ctx.Value(contextKey{}).(map[string]string); ok {
		data = meta
	} else {
		data = make(map[string]string)
	}

	for i := 0; i < len(kv); i += 2 {
		k := kv[i]
		v := kv[i+1]

		data[k] = v
	}

	return context.WithValue(ctx, contextKey{}, data)
}

// Wrap wraps an error with the metadata stored in the context using `WithMeta`.
func Wrap(err error, ctx context.Context) error {
	meta, ok := ctx.Value(contextKey{}).(map[string]string)
	if !ok {
		return err
	}

	return &withContext{err, meta}
}

// Unwrap pulls out any contextual metadata stored within an error.
func Unwrap(err error) map[string]string {
	values := map[string]string{}

	for err != nil {
		if f, ok := err.(*withContext); ok {
			if m := f.meta; m != nil {
				for k, v := range m {
					values[k] = v
				}
			}
		}

		err = errors.Unwrap(err)
	}

	if len(values) == 0 {
		return nil
	}

	return values
}
