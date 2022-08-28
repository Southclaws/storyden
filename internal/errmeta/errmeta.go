// Package errmeta allows storing metadata with errors.
package errmeta

import "errors"

// errmeta implements the Go error type and supports metadata that can easily be
// logged or sent as a response to clients.
type withMetadata struct {
	// the wrapped error value, either a standard library primitive or any other
	// error type from the ecosystem of error libraries.
	underlying error

	// a key-value pair much like context.valueCtx for storing any metadata.
	data map[string]string
}

// Implements all the interfaces for compatibility with the errors ecosystem.

func (e *withMetadata) Error() string {
	return e.underlying.Error()
}

func (e *withMetadata) ErrorData() map[string]string { return e.data }
func (e *withMetadata) Cause() error                 { return e.underlying }
func (e *withMetadata) Unwrap() error                { return e.underlying }
func (e *withMetadata) String() string               { return e.Error() }

// Wrap wraps an error along with a set of key-value pairs useful for describing
// the error in a structured way instead of with an unstructured string literal.
func Wrap(parent error, kv ...string) error {
	if parent == nil {
		return nil
	}

	if len(kv)%2 != 0 {
		panic("odd number of key-value pair arguments")
	}

	data := map[string]string{}

	for i := 0; i < len(kv); i += 2 {
		k := kv[i]
		v := kv[i+1]

		data[k] = v
	}

	return &withMetadata{
		underlying: parent,
		data:       data,
	}
}

// Metadata extracts any previously stored metadata from an error. If there was
// no metadata found then the return value is nil.
func Metadata(err error) map[string]string {
	values := map[string]string{}

	for err != nil {
		if f, ok := err.(interface {
			Metadata() map[string]string
		}); ok {
			if m := f.Metadata(); m != nil {
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
