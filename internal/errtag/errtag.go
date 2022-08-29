// Package errtag facilitates tagging error chains with distinct categories. The
// whole point of tagging errors with categories is to facilitate easier casting
// from errors to HTTP status codes at the transport layer. It also means you
// don't need to use explicit, manually defined custom errors for common things
// such as `sql.ErrNoRows`. You can decorate an error with metadata and tag it
// as a "not found" kind of error and at the transport layer, handle all errors
// with a unified error handler that can automatically set the HTTP status code.
//
// An error tag is any type which satisfies the interface `Tag() string`.
// Included in the library is a set of commonly used kinds of problem that can
// occur in most applications. These are based on gRPC status codes.
//
package errtag

import "errors"

type withTag struct {
	underlying error
	tag        errorTag
}

type errorTag interface {
	Tag() string
}

// Implements all the interfaces for compatibility with the errors ecosystem.

func (e *withTag) Error() string  { return e.underlying.Error() }
func (e *withTag) Cause() error   { return e.underlying }
func (e *withTag) Unwrap() error  { return e.underlying }
func (e *withTag) String() string { return e.Error() }

// Wrap wraps an error and gives it a distinct tag.
func Wrap(parent error, et errorTag) error {
	if parent == nil {
		return nil
	}

	if et == nil {
		return parent
	}

	return &withTag{
		underlying: parent,
		tag:        et,
	}
}

// Tag extracts the error tag of an error chain. If there's no tag, returns nil.
func Tag(err error) errorTag {
	for err != nil {
		if f, ok := err.(*withTag); ok {
			return f.tag
		}

		err = errors.Unwrap(err)
	}

	return nil
}

// Common kinds of error:

type Internal struct{}               // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
func (Internal) Tag() string         { return "INTERNAL" } //
type Cancelled struct{}              // The operation was cancelled, typically by the caller.
func (Cancelled) Tag() string        { return "CANCELLED" } //
type InvalidArgument struct{}        // The client specified an invalid argument.
func (InvalidArgument) Tag() string  { return "INVALID_ARGUMENT" } //
type NotFound struct{}               // Some requested entity was not found.
func (NotFound) Tag() string         { return "NOT_FOUND" } //
type AlreadyExists struct{}          // The entity that a client attempted to create already exists.
func (AlreadyExists) Tag() string    { return "ALREADY_EXISTS" } //
type PermissionDenied struct{}       // The caller does not have permission to execute the specified operation.
func (PermissionDenied) Tag() string { return "PERMISSION_DENIED" } //
type Unauthenticated struct{}        // The request does not have valid authentication credentials for the operation.
func (Unauthenticated) Tag() string  { return "UNAUTHENTICATED" } //
