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

type OK struct{}       // Not an error; returned on success.
func (OK) Tag() string { return "OK" }

type CANCELLED struct{}       // The operation was cancelled, typically by the caller.
func (CANCELLED) Tag() string { return "CANCELLED" }

type UNKNOWN struct{}       // Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
func (UNKNOWN) Tag() string { return "UNKNOWN" }

type INVALID_ARGUMENT struct{}       // The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
func (INVALID_ARGUMENT) Tag() string { return "INVALID_ARGUMENT" }

type DEADLINE_EXCEEDED struct{}       // The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long
func (DEADLINE_EXCEEDED) Tag() string { return "DEADLINE_EXCEEDED" }

type NOT_FOUND struct{}       // Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used.
func (NOT_FOUND) Tag() string { return "NOT_FOUND" }

type ALREADY_EXISTS struct{}       // The entity that a client attempted to create (e.g., file or directory) already exists.
func (ALREADY_EXISTS) Tag() string { return "ALREADY_EXISTS" }

type PERMISSION_DENIED struct{}       // The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some resource (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.
func (PERMISSION_DENIED) Tag() string { return "PERMISSION_DENIED" }

type RESOURCE_EXHAUSTED struct{}       // Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
func (RESOURCE_EXHAUSTED) Tag() string { return "RESOURCE_EXHAUSTED" }

type FAILED_PRECONDITION struct{}       // The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
func (FAILED_PRECONDITION) Tag() string { return "FAILED_PRECONDITION" }

type ABORTED struct{}       // The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
func (ABORTED) Tag() string { return "ABORTED" }

type OUT_OF_RANGE struct{}       // The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.
func (OUT_OF_RANGE) Tag() string { return "OUT_OF_RANGE" }

type UNIMPLEMENTED struct{}       // The operation is not implemented or is not supported/enabled in this service.
func (UNIMPLEMENTED) Tag() string { return "UNIMPLEMENTED" }

type INTERNAL struct{}       // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
func (INTERNAL) Tag() string { return "INTERNAL" }

type UNAVAILABLE struct{}       // The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.
func (UNAVAILABLE) Tag() string { return "UNAVAILABLE" }

type DATA_LOSS struct{}       // Unrecoverable data loss or corruption.
func (DATA_LOSS) Tag() string { return "DATA_LOSS" }

type UNAUTHENTICATED struct{}       // The request does not have valid authentication credentials for the operation.
func (UNAUTHENTICATED) Tag() string { return "UNAUTHENTICATED" }
