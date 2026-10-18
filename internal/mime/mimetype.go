package mime

import (
	"bytes"
	"io"
	"mime"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

const DefaultMIME = "application/octet-stream"

// Type represents a MIME type for asset file format identification.
type Type struct {
	mt string
}

func (t Type) String() string {
	return t.mt
}

func New(s string) Type {
	m := mimetype.Lookup(s)
	if m != nil {
		return Type{mt: m.String()}
	}

	mediaType, _, err := mime.ParseMediaType(s)
	if err != nil || !strings.Contains(mediaType, "/") {
		return Type{mt: DefaultMIME}
	}
	return Type{mt: mediaType}
}

func Detect(input io.Reader) (*Type, io.Reader, error) {
	header := bytes.NewBuffer(nil)

	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))
	if err != nil {
		return nil, nil, err
	}

	recycled := io.MultiReader(header, input)

	mt := &Type{
		mt: mtype.String(),
	}

	return mt, recycled, err
}
