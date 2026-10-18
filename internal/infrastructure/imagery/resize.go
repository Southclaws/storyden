package imagery

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

type ResizeMode uint8

const (
	FitMode ResizeMode = iota
	ThumbnailMode
	FillMode
)

func Resize(src image.Image, width, height int, square bool) *image.NRGBA {
	mode := FitMode
	if square {
		mode = ThumbnailMode
	}

	return resize(src, width, height, mode)
}

func Fill(src image.Image, width, height int) *image.NRGBA {
	return resize(src, width, height, FillMode)
}

func resize(src image.Image, width, height int, mode ResizeMode) *image.NRGBA {
	bounds := src.Bounds()
	if bounds.Empty() || width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	crop, destination := resizeGeometry(bounds, width, height, mode)
	dst := image.NewNRGBA(destination)
	draw.CatmullRom.Scale(dst, destination, src, crop, draw.Src, nil)
	return dst
}

func resizeGeometry(bounds image.Rectangle, width, height int, mode ResizeMode) (image.Rectangle, image.Rectangle) {
	if mode == FitMode {
		scale := min(1, float64(width)/float64(bounds.Dx()), float64(height)/float64(bounds.Dy()))
		width = max(1, int(float64(bounds.Dx())*scale))
		height = max(1, int(float64(bounds.Dy())*scale))

		return bounds, image.Rect(0, 0, width, height)
	}

	if mode == ThumbnailMode {
		size := min(width, height, bounds.Dx(), bounds.Dy())
		width, height = size, size
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

	return crop, image.Rect(0, 0, width, height)
}
