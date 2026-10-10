package avatar

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
)

type errReadStorer struct{}

func (errReadStorer) Exists(context.Context, string) (bool, error) { return false, nil }
func (errReadStorer) Read(context.Context, string) (io.Reader, int64, error) {
	return nil, 0, errors.New("not stored")
}
func (errReadStorer) Write(context.Context, string, io.Reader, int64) error { return nil }
func (errReadStorer) Delete(context.Context, string) error                  { return nil }
func (errReadStorer) List(context.Context, string) ([]string, error)        { return nil, nil }

type stubGenerator struct{}

func (stubGenerator) Generate(context.Context, string) (image.Image, error) {
	return image.NewNRGBA(image.Rect(0, 0, 2, 2)), nil
}

func TestGetFallbackStreamsGeneratedAvatar(t *testing.T) {
	t.Parallel()

	s := &service{generator: stubGenerator{}, storage: errReadStorer{}}

	r, _, err := s.Get(context.Background(), account.AccountID{})
	require.NoError(t, err)

	// before the fix this returned io.ErrClosedPipe instead of the png bytes
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	_, err = png.Decode(bytes.NewReader(data))
	require.NoError(t, err)
}