package errctx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMeta(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value"}, data)
}

func TestWithMetaOverwrite(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")
	ctx = WithMeta(ctx, "key", "value2")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value2"}, data)
}

func TestWithMetaNested(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")
	ctx = WithMeta(ctx, "key", "value2")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = WithMeta(ctx, "key", "value3")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]string{"key": "value3"}, data)
}

func TestWithMetaNestedManyKeys(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = WithMeta(ctx, "key2", "value2")
	ctx = WithMeta(ctx, "key3", "value3", "key4", "value4")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}, data)
}

func TestWithMetaEmpty(t *testing.T) {
	err := errors.New("a problem")
	data := Unwrap(err)

	assert.Nil(t, data)
}
