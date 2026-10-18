package imagery

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResizePreservesAspectRatio(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		source image.Rectangle
		width  int
		height int
		want   image.Rectangle
	}{
		{
			name:   "wide bounds",
			source: image.Rect(0, 0, 800, 400),
			width:  300,
			height: 100,
			want:   image.Rect(0, 0, 200, 100),
		},
		{
			name:   "round down",
			source: image.Rect(0, 0, 601, 400),
			width:  320,
			height: 320,
			want:   image.Rect(0, 0, 320, 212),
		},
		{name: "thin", source: image.Rect(0, 0, 1, 1000), width: 320, height: 320, want: image.Rect(0, 0, 1, 320)},
		{
			name:   "offset",
			source: image.Rect(10, 20, 110, 70),
			width:  320,
			height: 320,
			want:   image.Rect(0, 0, 100, 50),
		},
		{name: "empty source", width: 320, height: 320},
		{name: "empty target", source: image.Rect(0, 0, 100, 50)},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := image.NewNRGBA(test.source)
			result := Resize(source, test.width, test.height, false)
			require.Equal(t, test.want, result.Bounds())
		})
	}
}

func TestFillCropsCenterToFixedDimensions(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		source image.Rectangle
		center image.Rectangle
		width  int
		height int
	}{
		{
			name:   "landscape",
			source: image.Rect(10, 20, 310, 120),
			center: image.Rect(60, 20, 260, 120),
			width:  100,
			height: 50,
		},
		{
			name:   "portrait",
			source: image.Rect(10, 20, 110, 320),
			center: image.Rect(10, 145, 110, 195),
			width:  100,
			height: 50,
		},
		{
			name:   "upscale icon",
			source: image.Rect(0, 0, 2, 1),
			center: image.Rect(0, 0, 2, 1),
			width:  32,
			height: 32,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := image.NewNRGBA(test.source)
			for y := test.source.Min.Y; y < test.source.Max.Y; y++ {
				for x := test.source.Min.X; x < test.source.Max.X; x++ {
					pixel := color.NRGBA{R: 255, A: 255}
					if image.Pt(x, y).In(test.center) {
						pixel = color.NRGBA{G: 255, A: 255}
					}
					source.SetNRGBA(x, y, pixel)
				}
			}

			result := Fill(source, test.width, test.height)
			require.Equal(t, image.Rect(0, 0, test.width, test.height), result.Bounds())
			require.Equal(t, color.NRGBA{G: 255, A: 255}, result.NRGBAAt(test.width/2, test.height/2))
			require.Equal(t, color.NRGBA{G: 255, A: 255}, result.NRGBAAt(test.width/4, test.height/4))
		})
	}
}

func TestResizePreservesTransparentEdges(t *testing.T) {
	t.Parallel()

	source := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for y := 25; y < 75; y++ {
		for x := 25; x < 75; x++ {
			source.SetNRGBA(x, y, color.NRGBA{R: 255, A: 128})
		}
	}

	result := Resize(source, 50, 50, false)
	require.Zero(t, result.NRGBAAt(0, 0).A)
	require.Equal(t, color.NRGBA{R: 255, A: 128}, result.NRGBAAt(25, 25))
	edge := result.NRGBAAt(12, 25)
	require.Positive(t, edge.A)
	require.Less(t, edge.A, uint8(128))
	require.Equal(t, uint8(255), edge.R)
}

func TestEagerVariantMemoryAdmission(t *testing.T) {
	t.Parallel()

	source := &Source{
		Config: image.Config{Width: 4000, Height: 3000},
		Size:   8 << 20,
	}
	portrait := &Source{
		Config: image.Config{Width: 3000, Height: 4000},
		Size:   8 << 20,
	}
	thin := &Source{
		Config: image.Config{Width: 1, Height: 12000000},
		Size:   8 << 20,
	}

	const budget int64 = 256 << 20
	landscapeEstimate, err := Estimate(source, 320, 320, FitMode)
	require.NoError(t, err)
	require.Greater(t, landscapeEstimate.Total(), budget)

	smaller := &Source{
		Config: image.Config{Width: 1600, Height: 800},
		Size:   1 << 20,
	}
	smallEstimate, err := Estimate(smaller, 320, 320, FitMode)
	require.NoError(t, err)
	require.Less(t, smallEstimate.Total(), budget)

	portraitEstimate, err := Estimate(portrait, 320, 320, FitMode)
	require.NoError(t, err)
	require.Equal(t, landscapeEstimate.Total(), portraitEstimate.Total())

	thinEstimate, err := Estimate(thin, 320, 320, FitMode)
	require.NoError(t, err)
	require.Greater(t, thinEstimate.Total(), budget)
}

func TestEstimateCoversMeasuredMemory(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		width    int
		height   int
		encoded  int64
		target   int
		measured int64
	}{
		{
			name:     "transparent maximum with sequential variants",
			width:    4000,
			height:   3000,
			encoded:  300706,
			target:   1920,
			measured: 600 << 20,
		},
		{
			name:     "webp medium",
			width:    1024,
			height:   670,
			encoded:  74498,
			target:   960,
			measured: 90 << 20,
		},
		{
			name:     "thin image contribution tables",
			width:    1,
			height:   1000000,
			encoded:  28000,
			target:   320,
			measured: 100 << 20,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := &Source{
				Config: image.Config{Width: test.width, Height: test.height},
				Size:   test.encoded,
			}

			estimate, err := Estimate(source, test.target, test.target, FitMode)
			require.NoError(t, err)
			require.Greater(t, estimate.Total(), test.measured)
			require.Less(t, estimate.Total(), int64(768<<20))
		})
	}
}
