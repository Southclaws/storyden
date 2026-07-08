// package imagesafe decodes images with a pixel budget to prevent decompression bombs
package imagesafe

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

// MaxPixels caps decoded dimensions so a tiny file cannot allocate huge memory
const MaxPixels = 64 * 1000 * 1000

// maxEncodedBytes bounds how much encoded data we buffer before decoding
const maxEncodedBytes = 64 << 20

var (
	ErrImageTooLarge   = errors.New("image dimensions exceed the allowed size")
	ErrEncodedTooLarge = errors.New("encoded image exceeds the allowed size")
)

// Decode reads an image but rejects oversized dimensions before allocating pixels
func Decode(r io.Reader) (image.Image, string, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxEncodedBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxEncodedBytes {
		return nil, "", ErrEncodedTooLarge
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	// non-positive dimensions can appear from integer wrap on 32-bit and must not slip past the pixel budget
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > MaxPixels {
		return nil, format, ErrImageTooLarge
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, format, err
	}
	return img, format, nil
}
