package plugin_reader

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
)

type Reader struct {
	db    *ent.Client
	store object.Storer
	run   wrun.Runner
}

func New(
	db *ent.Client,
	store object.Storer,
	run wrun.Runner,
) *Reader {
	return &Reader{
		db:    db,
		store: store,
		run:   run,
	}
}

func (r *Reader) Get(ctx context.Context, id plugin.ID) (*plugin.Record, error) {
	record, err := r.db.Plugin.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := plugin.MapRecord(record)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	exists, err := r.store.Exists(ctx, rec.FilePath)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		rec.Status.ActiveState = plugin.ActiveStateError
		rec.Status.StatusMessage = "Plugin file not found in storage"
		rec.Status.Details = map[string]any{
			"expected_file_path": rec.FilePath,
			"current_file_paths": []string{}, // No files found
		}
	}

	return rec, nil
}

func (r *Reader) List(ctx context.Context) ([]*plugin.Record, error) {
	rs, err := r.db.Plugin.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Instead of erroring on a single record, we should collect all of
	// them anyway, and mark the errored ones as errored.
	records, err := dt.MapErr(rs, plugin.MapRecord)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	files, err := r.store.List(ctx, plugin.PluginDirectory)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	validated := dt.Map(records, func(r *plugin.Record) *plugin.Record {
		_, exists := lo.Find(files, func(f string) bool {
			return f == r.FilePath
		})
		if !exists {
			// NOTE: A bit of a mutative hack, these kinds of edge case error
			// states are not currently easier to represent in the data model.
			r.Status.ActiveState = plugin.ActiveStateError
			r.Status.StatusMessage = "Plugin file not found in storage"
			r.Status.Details = map[string]any{
				"expected_file_path": r.FilePath,
				"current_file_paths": files,
			}
		}

		return r
	})

	return validated, nil
}
