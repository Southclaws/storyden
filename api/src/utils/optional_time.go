package utils

import "4d63.com/optional"

type Zeroable interface {
	IsZero() bool
}

func OptionalZero[T Zeroable](t T) optional.Optional[T] {
	if t.IsZero() {
		return optional.Empty[T]()
	}

	return optional.Of(t)
}
