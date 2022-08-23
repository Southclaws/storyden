package utils

// ToMap is for use with samber/lo.Map
func ToMap[T any, R any](fn func(t T) R) func(t T, i int) R {
	return func(t T, i int) R { return fn(t) }
}

func Deref[T any](t *T, _ int) T {
	return *t
}

func Ref[T any](t T) *T {
	return &t
}
