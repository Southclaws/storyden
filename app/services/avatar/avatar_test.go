package avatar

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
)

type stubGenerator struct{}

func (stubGenerator) Generate(context.Context, string) (image.Image, error) {
	return image.NewNRGBA(image.Rect(0, 0, 2, 2)), nil
}

func TestGeneratedAvatarStreamsPNG(t *testing.T) {
	t.Parallel()

	s := &service{generator: stubGenerator{}}

	_, r, _, err := s.generated(context.Background(), account.AccountID{})
	require.NoError(t, err)

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	_, err = png.Decode(bytes.NewReader(data))
	require.NoError(t, err)
}
