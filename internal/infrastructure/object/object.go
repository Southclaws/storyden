package object

import (
	"context"
	"io"
)

type Storer interface {
	Exists(ctx context.Context, path string) (bool, error)
	Read(ctx context.Context, path string) (io.Reader, error)
	Write(ctx context.Context, path string, w io.Reader) error
}
