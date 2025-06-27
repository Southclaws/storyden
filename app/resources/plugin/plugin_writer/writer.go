package plugin_writer

import (
	"context"
	"io"
	"path/filepath"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
)

type Writer struct {
	store object.Storer
	run   wrun.Runner
}

func New(
	store object.Storer,
	run wrun.Runner,
) *Writer {
	return &Writer{
		store: store,
		run:   run,
	}
}

func (w *Writer) Create(ctx context.Context, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	pm, err := plugin.Binary(b).Validate(ctx, w.run)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	p := filepath.Join(plugin.PluginDirectory, pm.Name)

	err = w.store.Write(ctx, p, r, int64(len(b)))
	if err != nil {
		return err
	}

	return nil
}
