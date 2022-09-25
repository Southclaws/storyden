package object

import (
	"context"
	"io"
)

type Storer interface {
	Read(ctx context.Context, path string) (io.Reader, error)
	Write(ctx context.Context, path string, w io.Reader) error
}
