package utils

import (
	"4d63.com/optional"
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
