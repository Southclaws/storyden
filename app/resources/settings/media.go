package settings

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
)

var errInvalidMediaSettings = fault.New("invalid media settings")

type ImageLimits struct {
	MaxSize   int64
	MaxPixels int64
}

func (s MediaServiceSettings) Validate() error {
	if size, ok := s.ImageMaxSize.Get(); ok && size <= 0 {
		return fault.Wrap(errInvalidMediaSettings,
			fmsg.WithDesc(
				"nonpositive image maximum size",
				"The maximum image size must be greater than zero.",
			),
			ftag.With(ftag.InvalidArgument))
	}

	if pixels, ok := s.ImageMaxPixels.Get(); ok && pixels <= 0 {
		return fault.Wrap(errInvalidMediaSettings,
			fmsg.WithDesc(
				"nonpositive image maximum pixel count",
				"The maximum image pixel count must be greater than zero.",
			),
			ftag.With(ftag.InvalidArgument))
	}

	return nil
}

func (s *Settings) ImageLimits() (ImageLimits, error) {
	media := s.Services.OrZero().Media.OrZero()
	if err := media.Validate(); err != nil {
		return ImageLimits{}, err
	}

	size, sizeOK := media.ImageMaxSize.Get()
	pixels, pixelsOK := media.ImageMaxPixels.Get()
	if !sizeOK || !pixelsOK {
		return ImageLimits{}, fault.Wrap(errInvalidMediaSettings, fmsg.WithDesc(
			"missing image upload limits",
			"Image upload limits are not configured. Please contact an administrator.",
		))
	}

	return ImageLimits{MaxSize: size, MaxPixels: pixels}, nil
}
