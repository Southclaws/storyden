package imagery

import (
	"encoding/binary"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/bep/imagemeta"
	"github.com/gen2brain/webp"
	xwebp "golang.org/x/image/webp"
)

var (
	ErrUnsupported = errors.New("unsupported or animated image")
	ErrTooLarge    = errors.New("image exceeds processing limits")
)

type Source struct {
	File   *os.File
	Config image.Config
	Format string
	Size   int64
}

func Open(r io.Reader, maxBytes, maxPixels int64) (*Source, error) {
	f, err := os.CreateTemp("", "storyden-image-*")
	if err != nil {
		return nil, err
	}
	s := &Source{File: f}
	ok := false
	defer func() {
		if !ok {
			s.Close()
		}
	}()
	s.Size, err = io.Copy(f, io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if s.Size > maxBytes {
		return nil, ErrTooLarge
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	header := make([]byte, 12)
	_, _ = io.ReadFull(f, header)
	_, _ = f.Seek(0, io.SeekStart)
	switch {
	case string(header[:8]) == "\x89PNG\r\n\x1a\n":
		s.Format = "png"
		s.Config, err = png.DecodeConfig(f)
	case header[0] == 0xff && header[1] == 0xd8:
		s.Format = "jpeg"
		s.Config, err = jpeg.DecodeConfig(f)
	case string(header[:4]) == "RIFF" && string(header[8:]) == "WEBP":
		s.Format = "webp"
		s.Config, err = xwebp.DecodeConfig(f)
	default:
		return nil, ErrUnsupported
	}
	if err != nil {
		return nil, errors.Join(ErrUnsupported, err)
	}
	w, h := int64(s.Config.Width), int64(s.Config.Height)
	if w <= 0 || h <= 0 || w > maxPixels/h {
		return nil, ErrTooLarge
	}
	if err = s.checkAnimation(); err != nil {
		return nil, err
	}
	ok = true
	return s, nil
}

func (s *Source) checkAnimation() error {
	var offset int64
	if s.Format == "png" {
		offset = 8
	} else if s.Format == "webp" {
		offset = 12
	} else {
		return nil
	}
	for offset+8 <= s.Size {
		if _, err := s.File.Seek(offset, io.SeekStart); err != nil {
			return err
		}
		var h [8]byte
		if _, err := io.ReadFull(s.File, h[:]); err != nil {
			return err
		}
		var length int64
		var kind string
		if s.Format == "png" {
			length = int64(binary.BigEndian.Uint32(h[:4]))
			kind = string(h[4:])
			if kind == "acTL" {
				return ErrUnsupported
			}
			offset += 12 + length
		} else {
			kind = string(h[:4])
			length = int64(binary.LittleEndian.Uint32(h[4:]))
			if kind == "ANIM" || kind == "ANMF" {
				return ErrUnsupported
			}
			if kind == "VP8X" {
				var flags [1]byte
				if _, err := io.ReadFull(s.File, flags[:]); err != nil {
					return err
				}
				if flags[0]&2 != 0 {
					return ErrUnsupported
				}
			}
			offset += 8 + length + (length % 2)
		}
		if offset > s.Size {
			return ErrUnsupported
		}
	}
	return nil
}

func (s *Source) Decode() (image.Image, error) {
	if _, err := s.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if s.Format == "webp" {
		return webp.Decode(s.File, webp.Options{AutoRotate: true})
	}

	var decode func(io.Reader) (image.Image, error)
	var format imagemeta.ImageFormat
	switch s.Format {
	case "jpeg":
		decode = jpeg.Decode
		format = imagemeta.JPEG
	case "png":
		decode = png.Decode
		format = imagemeta.PNG
	default:
		return nil, ErrUnsupported
	}

	orientation := readOrientation(s.File, format)
	if _, err := s.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	decoded, err := decode(s.File)
	if err != nil {
		return nil, err
	}
	return orient(decoded, orientation), nil
}

func (s *Source) Close() {
	s.File.Close()
	os.Remove(s.File.Name())
}

func Estimate(s *Source, width, height int) int64 {
	// Account for orientation/conversion, one output, codec scratch and encoded copies.
	const sourcePixelBytes = 16
	const outputPixelBytes = 24
	const codecOverhead = 16 << 20
	const resizePixelBytes = 4 * 8

	sourceBytes := int64(s.Config.Width) * int64(s.Config.Height) * sourcePixelBytes
	outputBytes := int64(width) * int64(height) * outputPixelBytes
	encodedBytes := s.Size * 2

	// Bound Catmull-Rom's float64 workspace for fit/fill and either EXIF orientation.
	sourceWidth, sourceHeight := int64(s.Config.Width), int64(s.Config.Height)
	resizePixels := max(
		sourceWidth,
		sourceHeight,
		min(int64(width)*sourceHeight, int64(height)*sourceWidth),
		min(int64(width)*sourceWidth, int64(height)*sourceHeight),
	)
	resizeBytes := (resizePixels + int64(width)) * resizePixelBytes

	return sourceBytes + outputBytes + encodedBytes + codecOverhead + resizeBytes
}

func Encode(w io.Writer, img image.Image, pngOutput bool) error {
	if pngOutput {
		return png.Encode(w, img)
	}
	return webp.Encode(w, img, webp.Options{Quality: 80, Method: 3})
}
