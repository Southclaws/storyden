package imagery

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/webp"
)

func TestResizeAndEncode(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		sourceWidth, sourceHeight     int
		targetSize                    int
		expectedWidth, expectedHeight int
		square                        bool
	}{
		{"landscape", 1600, 800, 320, 320, 160, false},
		{"portrait", 800, 1600, 320, 160, 320, false},
		{"square", 800, 800, 320, 320, 320, false},
		{"small", 100, 50, 320, 100, 50, false},
		{"thumbnail", 1600, 800, 128, 128, 128, true},
		{"small_thumbnail", 80, 40, 128, 40, 40, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := image.NewNRGBA(image.Rect(0, 0, tc.sourceWidth, tc.sourceHeight))
			for y := 0; y < tc.sourceHeight; y++ {
				for x := 0; x < tc.sourceWidth; x++ {
					src.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8(y), 100, 128})
				}
			}
			resized := Resize(src, tc.targetSize, tc.targetSize, tc.square)
			require.Equal(t, image.Rect(0, 0, tc.expectedWidth, tc.expectedHeight), resized.Bounds())
			var out bytes.Buffer
			require.NoError(t, Encode(&out, resized, false))
			decoded, err := webp.Decode(bytes.NewReader(out.Bytes()))
			require.NoError(t, err)
			require.Equal(t, resized.Bounds(), decoded.Bounds())
			_, _, _, alpha := decoded.At(0, 0).RGBA()
			require.InDelta(t, 128*257, int(alpha), 257)
		})
	}
}

func TestThumbnailCropsCenter(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 300, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 300; x++ {
			c := color.NRGBA{R: 255, A: 255}
			if x >= 100 && x < 200 {
				c = color.NRGBA{G: 255, A: 255}
			}
			src.SetNRGBA(x, y, c)
		}
	}
	out := Resize(src, 100, 100, true)
	require.Equal(t, color.NRGBA{G: 255, A: 255}, out.NRGBAAt(50, 50))
	require.Equal(t, color.NRGBA{G: 255, A: 255}, out.NRGBAAt(1, 50))
}

func TestLimitsAndUnsupportedImages(t *testing.T) {
	var encoded bytes.Buffer
	require.NoError(t, png.Encode(&encoded, image.NewNRGBA(image.Rect(0, 0, 20, 10))))
	for _, tc := range []struct {
		name                string
		data                []byte
		maxBytes, maxPixels int64
		want                error
	}{
		{"encoded_limit", encoded.Bytes(), 8, 1000, ErrTooLarge},
		{"pixel_limit", encoded.Bytes(), 1 << 20, 100, ErrTooLarge},
		{"document", []byte("%PDF-1.7 document"), 1000, 1000, ErrUnsupported},
		{"zip", []byte("PK\x03\x04 archive"), 1000, 1000, ErrUnsupported},
		{"svg", []byte("<svg viewBox='0 0 10 10'/>"), 1000, 1000, ErrUnsupported},
		{"gif", []byte("GIF89a..."), 1000, 1000, ErrUnsupported},
		{"invalid", []byte("\xff\xd8 malformed JPEG"), 1000, 1000, ErrUnsupported},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Open(bytes.NewReader(tc.data), tc.maxBytes, tc.maxPixels)
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func TestAnimatedPNGIsSkippedBeforeDecode(t *testing.T) {
	var pngData bytes.Buffer
	require.NoError(t, png.Encode(&pngData, image.NewNRGBA(image.Rect(0, 0, 20, 10))))
	var chunk bytes.Buffer
	binary.Write(&chunk, binary.BigEndian, uint32(8))
	chunk.WriteString("acTL")
	binary.Write(&chunk, binary.BigEndian, uint32(2))
	binary.Write(&chunk, binary.BigEndian, uint32(0))
	binary.Write(&chunk, binary.BigEndian, crc32.ChecksumIEEE(chunk.Bytes()[4:]))
	data := append([]byte{}, pngData.Bytes()[:33]...)
	data = append(data, chunk.Bytes()...)
	data = append(data, pngData.Bytes()[33:]...)
	_, err := Open(bytes.NewReader(data), 1<<20, 1000)
	require.ErrorIs(t, err, ErrUnsupported)
}

func TestEXIFOrientation(t *testing.T) {
	var data bytes.Buffer
	require.NoError(t, jpeg.Encode(&data, image.NewNRGBA(image.Rect(0, 0, 40, 20)), nil))
	exif := []byte{
		'E', 'x', 'i', 'f', 0, 0,
		'I', 'I', 42, 0, 8, 0, 0, 0,
		1, 0,
		0x12, 1, 3, 0, 1, 0, 0, 0, 6, 0, 0, 0,
		0, 0, 0, 0,
	}
	oriented := append([]byte{}, data.Bytes()[:2]...)
	oriented = append(oriented, 0xff, 0xe1, byte((len(exif)+2)>>8), byte(len(exif)+2))
	oriented = append(oriented, exif...)
	oriented = append(oriented, data.Bytes()[2:]...)
	src, err := Open(bytes.NewReader(oriented), 1<<20, 1000)
	require.NoError(t, err)
	defer src.Close()
	img, err := src.Decode()
	require.NoError(t, err)
	require.Equal(t, image.Rect(0, 0, 20, 40), img.Bounds())
}

func TestPNGEXIFOrientation(t *testing.T) {
	t.Parallel()

	var encoded bytes.Buffer
	require.NoError(t, png.Encode(&encoded, image.NewNRGBA(image.Rect(0, 0, 40, 20))))
	for _, test := range []struct {
		name   string
		exif   []byte
		bounds image.Rectangle
	}{
		{
			name: "rotate clockwise",
			exif: []byte{
				'I', 'I', 42, 0, 8, 0, 0, 0,
				1, 0,
				0x12, 1, 3, 0, 1, 0, 0, 0, 6, 0, 0, 0,
				0, 0, 0, 0,
			},
			bounds: image.Rect(0, 0, 20, 40),
		},
		{
			name:   "invalid metadata preserves pixels",
			exif:   []byte("invalid EXIF"),
			bounds: image.Rect(0, 0, 40, 20),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var chunk bytes.Buffer
			require.NoError(t, binary.Write(&chunk, binary.BigEndian, uint32(len(test.exif))))
			chunk.WriteString("eXIf")
			chunk.Write(test.exif)
			require.NoError(t, binary.Write(&chunk, binary.BigEndian, crc32.ChecksumIEEE(chunk.Bytes()[4:])))
			data := append([]byte{}, encoded.Bytes()[:33]...)
			data = append(data, chunk.Bytes()...)
			data = append(data, encoded.Bytes()[33:]...)

			source, err := Open(bytes.NewReader(data), 1<<20, 1000)
			require.NoError(t, err)
			defer source.Close()
			decoded, err := source.Decode()
			require.NoError(t, err)
			require.Equal(t, test.bounds, decoded.Bounds())
		})
	}
}

func BenchmarkImageVariants(b *testing.B) {
	src := image.NewNRGBA(image.Rect(0, 0, 4000, 3000))
	for y := 0; y < 3000; y++ {
		for x := 0; x < 4000; x++ {
			src.SetNRGBA(x, y, color.NRGBA{uint8(x*17 + y*7), uint8(x + y*11), uint8(y * 5), 255})
		}
	}
	var data bytes.Buffer
	jpeg.Encode(&data, src, &jpeg.Options{Quality: 90})
	b.ReportAllocs()
	b.SetBytes(int64(data.Len()))
	b.ResetTimer()
	for range b.N {
		input, err := Open(bytes.NewReader(data.Bytes()), 32<<20, 12000000)
		if err != nil {
			b.Fatal(err)
		}
		img, err := input.Decode()
		if err != nil {
			b.Fatal(err)
		}
		for _, n := range []int{64, 128, 256, 320, 960, 1920} {
			resized := Resize(img, n, n, n < 320)
			var out bytes.Buffer
			if err := Encode(&out, resized, false); err != nil {
				b.Fatal(err)
			}
		}
		input.Close()
	}
}

func BenchmarkEagerImageVariants(b *testing.B) {
	img := image.NewNRGBA(image.Rect(0, 0, 4000, 3000))
	for y := 0; y < 3000; y++ {
		for x := 0; x < 4000; x++ {
			img.SetNRGBA(x, y, color.NRGBA{uint8(x*17 + y*7), uint8(x + y*11), uint8(y * 5), 255})
		}
	}
	var encoded bytes.Buffer
	if err := jpeg.Encode(&encoded, img, &jpeg.Options{Quality: 90}); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		source, err := Open(bytes.NewReader(encoded.Bytes()), 32<<20, 12000000)
		if err != nil {
			b.Fatal(err)
		}
		decoded, err := source.Decode()
		if err != nil {
			b.Fatal(err)
		}
		for _, size := range []int{64, 128, 256, 320} {
			var out bytes.Buffer
			if err = Encode(&out, Resize(decoded, size, size, size < 320), false); err != nil {
				b.Fatal(err)
			}
		}
		source.Close()
	}
}
