package utils

import (
	"4d63.com/optional"
	"github.com/rs/xid"
)

type Zeroable interface {
	IsZero() bool
}

func OptionalZero[T Zeroable](t T) optional.Optional[T] {
	if t.IsZero() {
		return optional.Empty[T]()
	}

	return optional.Of(t)
}

func OptionalPointer[T Zeroable](t T) optional.Optional[T] {
	return optional.Of(t)
}

func OptionalSlice[T any](t []T) optional.Optional[[]T] {
	if t == nil {
		return optional.Empty[[]T]()
	}

	return optional.Of(t)
}

func OptionalToPointer[T any](o optional.Optional[T]) *T {
	if v, ok := o.Get(); ok {
		return &v
	}

	return nil
}

func OptionalElse[T, R any](o optional.Optional[T], fn func(T) R) R {
	if v, ok := o.Get(); ok {
		r := fn(v)
		return r
	}

	return *new(R)
}

func OptionalElsePtr[T, R any](o optional.Optional[T], fn func(T) R) *R {
	if v, ok := o.Get(); ok {
		r := fn(v)
		return &r
	}

	return nil
}

// OptionalID is a special helper only used for tests where the seed data must
// set IDs for resources by calling `WithID`. In this case, the ID must be set
// using Ent's `SetNillableID` builder method. This function works by taking a
// a pointer to an ID and only setting it if the pointer is not nil. The problem
// with this is that all resource structs use IDs as values, not pointers so in
// order to derive a pointer that may or may not be nil, this function simply
// checks if the xid is considered "valid" and if it isn't, simply returns nil.
func OptionalID(id xid.ID) *xid.ID {
	if id.IsNil() {
		return nil
	}

	return &id
}
