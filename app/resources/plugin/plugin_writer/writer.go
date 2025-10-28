package plugin_writer

import (
	"bytes"
	"context"
	"path/filepath"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Writer struct {
	db    *ent.Client
	store object.Storer
}

func New(
	db *ent.Client,
	store object.Storer,
) *Writer {
	return &Writer{
		db:    db,
		store: store,
	}
}

func (w *Writer) Add(ctx context.Context, acc account.AccountID, pl *plugin.Validated) (*plugin.Available, error) {
	p := filepath.Join(plugin.PluginDirectory, string(pl.Metadata.ID))

	b := pl.Binary

	err := w.store.Write(ctx, p, bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := w.db.Plugin.Create().
		SetAccountID(xid.ID(acc)).
		SetConfig(map[string]any{}).
		SetManifest(pl.Metadata).
		SetActiveState(plugin.ActiveStateInactive.String()).
		SetActiveStateChangedAt(time.Now()).
		SetPath(p).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = fault.Wrap(err, ftag.With(ftag.AlreadyExists))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := plugin.MapRecord(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &plugin.Available{
		Record: *record,
		Loaded: pl,
	}, nil
}

func (w *Writer) Remove(ctx context.Context, plid plugin.ID) error {
	r, err := w.db.Plugin.Get(ctx, xid.ID(plid))
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err = w.store.Delete(ctx, r.Path); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err = w.db.Plugin.DeleteOne(r).Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
