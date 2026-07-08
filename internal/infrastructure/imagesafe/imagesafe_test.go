package imagesafe

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeAllowsNormalImage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 32, 32))))

	img, format, err := Decode(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 32, img.Bounds().Dx())
}

func TestDecodeRejectsOversizedDimensions(t *testing.T) {
	t.Parallel()

	// a valid png header claiming enormous dimensions must be rejected before the
	// pixel buffer is allocated, which would otherwise consume gigabytes
	side := 200000
	assert.Greater(t, int64(side)*int64(side), int64(MaxPixels))

	_, _, err := Decode(bytes.NewReader(pngHeader(side, side)))
	require.ErrorIs(t, err, ErrImageTooLarge)
}

// pngHeader builds a valid png signature and IHDR chunk so DecodeConfig can read
// the declared dimensions without any pixel data present
func pngHeader(w, h int) []byte {
	var b bytes.Buffer
	b.Write([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})

	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:], uint32(w))
	binary.BigEndian.PutUint32(ihdr[4:], uint32(h))
	ihdr[8] = 8 // bit depth
	ihdr[9] = 6 // colour type rgba

	binary.Write(&b, binary.BigEndian, uint32(len(ihdr)))
	chunk := append([]byte("IHDR"), ihdr...)
	b.Write(chunk)
	binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(chunk))

	return b.Bytes()
}
