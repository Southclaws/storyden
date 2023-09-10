package object

import (
	"context"
	"io"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

type Storer interface {
	Exists(ctx context.Context, path string) (bool, error)
	Read(ctx context.Context, path string) (io.Reader, int64, error)
	Write(ctx context.Context, path string, w io.Reader, size int64) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config) Storer {
		switch cfg.AssetStorageType {
		case "s3":
			return NewS3Storer(cfg)

		default:
			return NewLocalStorer(cfg)
		}
	})
}
