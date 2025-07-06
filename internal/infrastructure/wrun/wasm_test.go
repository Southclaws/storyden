package wrun

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_wazeroRunner_RunOnce(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	ctx := context.Background()

	runner := newWazeroRunner(ctx)

	f, err := os.ReadFile("./testdata/test1.wasm")
	r.NoError(err)

	want := []byte(`{"name":"test_plugin","version":"1.0"}`)

	b, err := runner.RunOnce(ctx, f, nil)
	r.NoError(err)
	a.Equal(want, b)
}

func Test_wazeroRunner_RunOnce2(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	ctx := context.Background()

	runner := newWazeroRunner(ctx)

	f, err := os.ReadFile("./testdata/test2.wasm")
	r.NoError(err)

	input := map[string]any{
		"event": "ThreadCreate",
		"id":    "123",
	}
	want := []byte(`{"event":"ThreadCreate","id":"123"}`)

	b, err := runner.RunOnce(ctx, f, input)
	r.NoError(err)
	a.Equal(want, b)
}
