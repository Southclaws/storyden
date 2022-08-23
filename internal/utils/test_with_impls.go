package utils

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

type ImplConstructor[T any] func() T

// TestAll runs each test for each implementation of some interface.
func TestAll[T any](
	t *testing.T,
	implementations []ImplConstructor[T],
	fn func(*testing.T, T),
) {
	lo.ForEach(implementations, func(cons ImplConstructor[T], _ int) {
		i := cons()

		name := reflect.TypeOf(i).Elem().Name()
		t.Run(name, func(t *testing.T) {
			fn(t, i)
		})
	})
}
