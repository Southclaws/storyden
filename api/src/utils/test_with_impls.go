package utils

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

// TestAll runs each test for each implementation of some interface.
func TestAll[T any](
	t *testing.T,
	implementations []T,
	fn func(*testing.T, T),
) {
	lo.ForEach(implementations, func(i T, _ int) {
		name := reflect.TypeOf(i).Elem().Name()
		t.Run(name, func(t *testing.T) {
			fn(t, i)
		})
	})
}

func TestAllNamed[T any](
	t *testing.T,
	name string,
	implementations []T,
	fn func(*testing.T, T),
) {
	t.Run(name, func(t *testing.T) {
		TestAll(t, implementations, fn)
	})
}
