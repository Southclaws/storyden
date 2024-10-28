package mime

import (
	"bytes"
	"io"

	"github.com/gabriel-vasile/mimetype"
)

const DefaultMIME = "application/octet-stream"

// Type represents a MIME type for asset file format identification.
type Type struct {
	mt mimetype.MIME
}

func (t Type) String() string {
	return t.mt.String()
}

func New(s string) Type {
	return Type{
		mt: *mimetype.Lookup(s),
	}
}

func Detect(input io.Reader) (*Type, io.Reader, error) {
	header := bytes.NewBuffer(nil)

	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))
	if err != nil {
		return nil, nil, err
	}

	recycled := io.MultiReader(header, input)

	mt := &Type{
		mt: *mtype,
	}

	return mt, recycled, err
}
