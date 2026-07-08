// package imagesafe decodes images with a pixel budget to prevent decompression bombs
package imagesafe

import (
	"bytes"
	"errors"
	"fmt"
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

var ErrImageTooLarge = errors.New("image dimensions exceed the allowed size")

// Decode reads an image but rejects oversized dimensions before allocating pixels
func Decode(r io.Reader) (image.Image, string, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxEncodedBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxEncodedBytes {
		return nil, "", fmt.Errorf("image exceeds %d bytes", int64(maxEncodedBytes))
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	if int64(cfg.Width)*int64(cfg.Height) > MaxPixels {
		return nil, format, ErrImageTooLarge
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, format, err
	}
	return img, format, nil
}
