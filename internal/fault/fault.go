// Package fault is a vendored version of fault to experiment with the API
// design in a real world application before moving it out to the public repo.
package fault

// fault implements the Go error type and supports metadata that can easily be
// logged or sent as a response to clients.
type fault struct {
	// the wrapped error value, either a standard library primitive or any other
	// error type from the ecosystem of error libraries.
	underlying error

	// location context of this particular error context so we don't need to
	// store a full stack trace of mostly useless info.
	location string

	// a key-value pair much like context.valueCtx for storing any metadata.
	data map[string]string
}

// Implements all the interfaces for compatibility with the errors ecosystem.

func (e *fault) Error() string {
	return e.underlying.Error()
}

func (e *fault) ErrorData() map[string]string { return e.data }
func (e *fault) Location() string             { return e.location }
func (e *fault) Cause() error                 { return e.underlying }
func (e *fault) Unwrap() error                { return e.underlying }
func (e *fault) String() string               { return e.Error() }

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

	return &fault{
		underlying: parent,
		data:       data,
		location:   getLocation(),
	}
}
