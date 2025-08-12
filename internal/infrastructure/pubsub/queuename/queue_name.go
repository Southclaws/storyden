package queuename

import "reflect"

func FromValue(zero any) string {
	t := reflect.TypeOf(zero)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	to := t.String()

	return to
}

func FromT[T any]() string {
	var zero T
	t := reflect.TypeOf(zero)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	to := t.String()

	return to
}
