package imagery

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrient(t *testing.T) {
	t.Parallel()

	source := image.NewNRGBA(image.Rect(10, 20, 13, 22))
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			source.SetNRGBA(10+x, 20+y, color.NRGBA{R: uint8(y*3 + x + 1), A: 255})
		}
	}

	for _, test := range []struct {
		name        string
		orientation uint16
		rows        [][]uint8
	}{
		{name: "normal", orientation: 1, rows: [][]uint8{{1, 2, 3}, {4, 5, 6}}},
		{name: "horizontal mirror", orientation: 2, rows: [][]uint8{{3, 2, 1}, {6, 5, 4}}},
		{name: "rotate 180", orientation: 3, rows: [][]uint8{{6, 5, 4}, {3, 2, 1}}},
		{name: "vertical mirror", orientation: 4, rows: [][]uint8{{4, 5, 6}, {1, 2, 3}}},
		{name: "transpose", orientation: 5, rows: [][]uint8{{1, 4}, {2, 5}, {3, 6}}},
		{name: "rotate clockwise", orientation: 6, rows: [][]uint8{{4, 1}, {5, 2}, {6, 3}}},
		{name: "transverse", orientation: 7, rows: [][]uint8{{6, 3}, {5, 2}, {4, 1}}},
		{name: "rotate counterclockwise", orientation: 8, rows: [][]uint8{{3, 6}, {2, 5}, {1, 4}}},
		{name: "invalid", orientation: 9, rows: [][]uint8{{1, 2, 3}, {4, 5, 6}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := orient(source, test.orientation)
			bounds := result.Bounds()
			require.Equal(t, len(test.rows), bounds.Dy())
			require.Equal(t, len(test.rows[0]), bounds.Dx())
			for y, row := range test.rows {
				for x, value := range row {
					pixel := color.NRGBAModel.Convert(result.At(bounds.Min.X+x, bounds.Min.Y+y))
					require.Equal(t, color.NRGBA{R: value, A: 255}, pixel, "pixel (%d, %d)", x, y)
				}
			}
		})
	}
}
