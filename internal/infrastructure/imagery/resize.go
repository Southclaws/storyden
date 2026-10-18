package imagery

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

func Resize(src image.Image, width, height int, square bool) *image.NRGBA {
	bounds := src.Bounds()
	if bounds.Empty() || width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	if square {
		size := min(width, height, bounds.Dx(), bounds.Dy())
		return Fill(src, size, size)
	}

	scale := min(1, float64(width)/float64(bounds.Dx()), float64(height)/float64(bounds.Dy()))
	width = max(1, int(float64(bounds.Dx())*scale))
	height = max(1, int(float64(bounds.Dy())*scale))

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Src, nil)
	return dst
}

func Fill(src image.Image, width, height int) *image.NRGBA {
	bounds := src.Bounds()
	if bounds.Empty() || width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	cropWidth, cropHeight := bounds.Dx(), bounds.Dy()
	sourceAspect := float64(cropWidth) / float64(cropHeight)
	targetAspect := float64(width) / float64(height)
	if sourceAspect > targetAspect {
		cropWidth = max(1, int(math.Round(float64(cropHeight)*targetAspect)))
	} else {
		cropHeight = max(1, int(math.Round(float64(cropWidth)/targetAspect)))
	}

	left := bounds.Min.X + (bounds.Dx()-cropWidth)/2
	top := bounds.Min.Y + (bounds.Dy()-cropHeight)/2
	crop := image.Rect(left, top, left+cropWidth, top+cropHeight)

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, crop, draw.Src, nil)
	return dst
}
