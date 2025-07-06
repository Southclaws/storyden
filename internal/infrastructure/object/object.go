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
	Delete(ctx context.Context, path string) error
	List(ctx context.Context, prefix string) ([]string, error)
}

func Build() fx.Option {
	return fx.Provide(func(ctx context.Context, cfg config.Config) (Storer, error) {
		switch cfg.AssetStorageType {
		case "s3":
			return NewS3Storer(ctx, cfg)

		default:
			return NewLocalStorer(cfg), nil
		}
	})
}
