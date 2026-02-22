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
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Writer struct {
	db             *ent.Client
	store          object.Storer
	pluginDataPath string
}

func New(
	cfg config.Config,
	db *ent.Client,
	store object.Storer,
) *Writer {
	return &Writer{
		db:             db,
		store:          store,
		pluginDataPath: cfg.PluginDataPath,
	}
}

func (r *Writer) FilePath(id plugin.InstallationID) string {
	return filepath.Join(plugin.PluginDirectory, id.String()+".zip")
}

func (w *Writer) Add(ctx context.Context, acc account.AccountID, pl *plugin.Validated) (*plugin.Available, error) {
	authSecret, err := plugin_auth.GenerateSecret()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := w.db.Plugin.Create().
		SetAccountID(xid.ID(acc)).
		SetConfig(map[string]any{}).
		SetManifest(pl.Metadata.ToMap()).
		SetSupervised(true).
		SetActiveState(plugin.ActiveStateInactive.String()).
		SetActiveStateChangedAt(time.Now()).
		SetAuthSecret(authSecret).
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

	p := w.FilePath(record.InstallationID)
	b := pl.Binary

	err = w.store.Write(ctx, p, bytes.NewReader(b), int64(len(b)))
	if err != nil {
		w.db.Plugin.DeleteOne(r).Exec(ctx)
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &plugin.Available{
		Record: *record,
		Loaded: pl,
	}, nil
}

func (w *Writer) Remove(ctx context.Context, plid plugin.InstallationID) error {
	r, err := w.db.Plugin.Get(ctx, xid.ID(plid))
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if r.Supervised {
		p := w.FilePath(plid)
		if err = w.store.Delete(ctx, p); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err = w.db.Plugin.DeleteOne(r).Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) AddExternal(
	ctx context.Context,
	acc account.AccountID,
	manifest rpc.Manifest,
) (*plugin.Record, string, error) {
	token, err := plugin_auth.GenerateExternalToken()
	if err != nil {
		return nil, "", err
	}

	now := time.Now()

	r, err := w.db.Plugin.Create().
		SetAccountID(xid.ID(acc)).
		SetConfig(map[string]any{}).
		SetManifest(manifest.ToMap()).
		SetSupervised(false).
		SetActiveState(plugin.ActiveStateActive.String()).
		SetActiveStateChangedAt(now).
		SetAuthSecret(token).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = fault.Wrap(err, ftag.With(ftag.AlreadyExists))
		}
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	record, err := plugin.MapRecord(r)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	return record, token, nil
}

func (w *Writer) SetActiveState(ctx context.Context, plid plugin.InstallationID, state plugin.ActiveState) error {
	_, err := w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetActiveState(state.String()).
		SetActiveStateChangedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) CycleAuthSecret(ctx context.Context, plid plugin.InstallationID) (string, error) {
	newSecret, err := plugin_auth.GenerateSecret()
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	_, err = w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetAuthSecret(newSecret).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return newSecret, nil
}

func (w *Writer) CycleExternalToken(ctx context.Context, plid plugin.InstallationID) (string, error) {
	newToken, err := plugin_auth.GenerateExternalToken()
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	_, err = w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetAuthSecret(newToken).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return newToken, nil
}

func (w *Writer) UpdateManifest(
	ctx context.Context,
	plid plugin.InstallationID,
	manifest rpc.Manifest,
) (*plugin.Record, error) {
	r, err := w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetManifest(manifest.ToMap()).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := plugin.MapRecord(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rec, nil
}

func (w *Writer) UpdatePackage(
	ctx context.Context,
	plid plugin.InstallationID,
	pl *plugin.Validated,
) (*plugin.Record, error) {
	existing, err := w.db.Plugin.Get(ctx, xid.ID(plid))
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !existing.Supervised {
		return nil, fault.Wrap(
			fault.New("cannot update package for external plugin"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	archivePath := w.FilePath(plid)
	b := pl.Binary
	if err := w.store.Write(ctx, archivePath, bytes.NewReader(b), int64(len(b))); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetManifest(pl.Metadata.ToMap()).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := plugin.MapRecord(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rec, nil
}

func (w *Writer) UpdateConfig(
	ctx context.Context,
	plid plugin.InstallationID,
	config map[string]any,
) (*plugin.Record, error) {
	r, err := w.db.Plugin.UpdateOneID(xid.ID(plid)).
		SetConfig(config).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := plugin.MapRecord(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rec, nil
}
