package imagery

import (
	"image"
	"io"

	"github.com/bep/imagemeta"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

func readOrientation(reader io.ReadSeeker, format imagemeta.ImageFormat) uint16 {
	orientation := uint16(1)
	_, err := imagemeta.Decode(imagemeta.Options{
		R:            reader,
		ImageFormat:  format,
		Sources:      imagemeta.EXIF,
		LimitTagSize: 2,
		ShouldHandleTag: func(tag imagemeta.TagInfo) bool {
			return tag.Tag == "Orientation" && tag.Namespace == "IFD0"
		},
		HandleTag: func(tag imagemeta.TagInfo) error {
			if value, ok := tag.Value.(uint16); ok {
				orientation = value
			}
			return imagemeta.ErrStopWalking
		},
	})
	if err != nil {
		// Invalid metadata must not prevent decoding otherwise valid pixels.
		return 1
	}
	return orientation
}

func orient(src image.Image, orientation uint16) image.Image {
	if orientation < 2 || orientation > 8 {
		return src
	}

	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	w, h := float64(width), float64(height)
	var transform f64.Aff3
	switch orientation {
	case 2:
		transform = f64.Aff3{-1, 0, w, 0, 1, 0}
	case 3:
		transform = f64.Aff3{-1, 0, w, 0, -1, h}
	case 4:
		transform = f64.Aff3{1, 0, 0, 0, -1, h}
	case 5:
		transform = f64.Aff3{0, 1, 0, 1, 0, 0}
	case 6:
		transform = f64.Aff3{0, -1, h, 1, 0, 0}
	case 7:
		transform = f64.Aff3{0, -1, h, -1, 0, w}
	case 8:
		transform = f64.Aff3{0, 1, 0, -1, 0, w}
	}
	if orientation >= 5 {
		width, height = height, width
	}

	// The EXIF transform uses coordinates relative to the source's top-left corner.
	transform[2] -= transform[0]*float64(bounds.Min.X) + transform[1]*float64(bounds.Min.Y)
	transform[5] -= transform[3]*float64(bounds.Min.X) + transform[4]*float64(bounds.Min.Y)
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Transform(dst, transform, src, bounds, draw.Src, nil)
	return dst
}
