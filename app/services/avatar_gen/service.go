package avatar_gen

import (
	"context"
	"image"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/mazznoer/colorgrad"
	"github.com/mazznoer/csscolorparser"
)

type AvatarGenerator interface {
	Generate(ctx context.Context, handle string) (image.Image, error)
}

func New() AvatarGenerator {
	return &service{}
}

type service struct{}

var start = csscolorparser.FromHsl(216.0, 0.1, 0.2, 1.0)

func (s *service) Generate(ctx context.Context, handle string) (image.Image, error) {
	hash := hashfunction(handle)

	c2 := csscolorparser.FromHsl(float64(hash), 0.69, 0.4, 1.0)

	grad, err := colorgrad.
		NewGradient().
		Colors(start, c2).
		Interpolation(colorgrad.InterpolationCatmullRom).
		Build()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	w := 128
	h := 128
	fw := float64(w)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			col := grad.At(float64(x+y) / (fw * 2))
			img.Set(x, y, col)
		}
	}

	return img, nil
}

func hashfunction(id string) uint16 {
	return dt.Reduce([]byte(id), func(r uint16, b byte) uint16 {
		s := uint16(b) * 42
		x := (r + 1) * s % 360
		return x
	}, 69)
}
