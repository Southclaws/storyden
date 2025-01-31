package plugin_reader

import (
	"context"
	"io"
	"path"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
)

type Reader struct {
	store object.Storer
	run   wrun.Runner
}

func New(
	store object.Storer,
	run wrun.Runner,
) *Reader {
	return &Reader{
		store: store,
		run:   run,
	}
}

func (r *Reader) List(ctx context.Context) ([]*plugin.Available, error) {
	files, err := r.store.List(ctx, plugin.PluginDirectory)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ps := dt.Map(files, func(s string) *plugin.Available {
		p := path.Join(plugin.PluginDirectory, s)

		br, _, err := r.store.Read(ctx, p)
		if err != nil {
			return &plugin.Available{Error: err}
		}

		b, err := io.ReadAll(br)
		if err != nil {
			return &plugin.Available{Error: err}
		}

		m, err := plugin.Binary(b).Validate(ctx, r.run)
		if err != nil {
			return &plugin.Available{Error: err}
		}

		return &plugin.Available{
			StoragePath: p,
			Loaded: &plugin.Package{
				Metadata: *m,
				Binary:   b,
			},
		}
	})

	return ps, nil
}
