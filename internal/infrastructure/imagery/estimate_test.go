package imagery

import (
	"image"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEstimateUsesTransformGeometry(t *testing.T) {
	source := &Source{
		Config: image.Config{Width: 3072, Height: 1536},
		Size:   1 << 20,
	}
	fitted, err := Estimate(source, 1920, 1920, FitMode)
	require.NoError(t, err)
	filled, err := Estimate(source, 1920, 1920, FillMode)
	require.NoError(t, err)
	require.Less(t, fitted.Output, filled.Output)

	small := &Source{
		Config: image.Config{Width: 80, Height: 40},
		Size:   1024,
	}
	thumbnail, err := Estimate(small, 256, 256, ThumbnailMode)
	require.NoError(t, err)
	icon, err := Estimate(small, 256, 256, FillMode)
	require.NoError(t, err)
	require.Less(t, thumbnail.Output, icon.Output)
	require.Less(t, thumbnail.Resize, icon.Resize)
}

func TestEstimateRejectsInvalidMemoryRequirements(t *testing.T) {
	for _, test := range []struct {
		name   string
		source *Source
		width  int
		height int
		mode   ResizeMode
	}{
		{name: "nil source", width: 320, height: 320},
		{name: "empty source", source: &Source{}, width: 320, height: 320},
		{
			name:   "negative output",
			source: &Source{Config: image.Config{Width: 100, Height: 100}},
			width:  -1,
			height: 320,
		},
		{
			name:   "source allocation overflow",
			source: &Source{Config: image.Config{Width: math.MaxInt, Height: math.MaxInt}},
			width:  320,
			height: 320,
		},
		{
			name:   "encoded allocation overflow",
			source: &Source{Config: image.Config{Width: 100, Height: 100}, Size: math.MaxInt64},
			width:  320,
			height: 320,
		},
		{
			name:   "output allocation overflow",
			source: &Source{Config: image.Config{Width: 100, Height: 100}},
			width:  math.MaxInt,
			height: math.MaxInt,
			mode:   FillMode,
		},
		{
			name:   "safety allowance overflow",
			source: &Source{Config: image.Config{Width: 100, Height: 100}, Size: math.MaxInt64 / 3},
			width:  320,
			height: 320,
		},
		{
			name:   "resampler contributions exceed int32 indexing",
			source: &Source{Config: image.Config{Width: 600000000, Height: 1}},
			width:  320,
			height: 320,
		},
		{
			name:   "resampler workspace exceeds int32 indexing",
			source: &Source{Config: image.Config{Width: 50000, Height: 50000}},
			width:  50000,
			height: 50000,
			mode:   FillMode,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Estimate(test.source, test.width, test.height, test.mode)
			require.ErrorIs(t, err, ErrTooLarge)
		})
	}
}
