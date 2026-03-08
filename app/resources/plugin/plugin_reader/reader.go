package plugin_reader

import (
	"context"
	"io"
	"path/filepath"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_plugin "github.com/Southclaws/storyden/internal/ent/plugin"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Reader struct {
	db             *ent.Client
	store          object.Storer
	pluginDataPath string
}

func New(
	cfg config.Config,
	db *ent.Client,
	store object.Storer,
) *Reader {
	return &Reader{
		db:             db,
		store:          store,
		pluginDataPath: cfg.PluginDataPath,
	}
}

func (r *Reader) FilePath(id plugin.InstallationID) string {
	return filepath.Join(plugin.PluginDirectory, id.String()+".zip")
}

func (r *Reader) Get(ctx context.Context, id plugin.InstallationID) (*plugin.Record, error) {
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

	if rec.Mode.Supervised() {
		filePath := r.FilePath(rec.InstallationID)
		exists, err := r.store.Exists(ctx, filePath)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		if !exists {
			rec.StatusMessage = "Plugin file not found in storage"
			rec.Details = map[string]any{
				"expected_file_path": filePath,
			}
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

	validated := dt.Map(records, func(rec *plugin.Record) *plugin.Record {
		if !rec.Mode.Supervised() {
			return rec
		}

		filePath := r.FilePath(rec.InstallationID)
		_, exists := lo.Find(files, func(f string) bool {
			path := filepath.Join(plugin.PluginDirectory, f)
			return path == filePath
		})
		if !exists {
			rec.StatusMessage = "Plugin file not found in storage"
			rec.Details = map[string]any{
				"expected_file_path": filePath,
				"current_file_paths": files,
			}
		}

		return rec
	})

	return validated, nil
}

func (r *Reader) LoadBinary(ctx context.Context, id plugin.InstallationID) ([]byte, error) {
	filePath := r.FilePath(id)

	reader, _, err := r.store.Read(ctx, filePath)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	bin, err := io.ReadAll(reader)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return bin, nil
}

func (r *Reader) GetAuthSecret(ctx context.Context, id plugin.InstallationID) (string, error) {
	record, err := r.db.Plugin.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return "", fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return record.AuthSecret, nil
}

func (r *Reader) GetByExternalToken(ctx context.Context, token string) (*plugin.Record, error) {
	record, err := r.db.Plugin.Query().
		Where(
			ent_plugin.SupervisedEQ(false),
			ent_plugin.AuthSecretEQ(token),
		).
		Only(ctx)
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

	return rec, nil
}

func (r *Reader) GetConfig(ctx context.Context, id plugin.InstallationID) (map[string]any, error) {
	record, err := r.db.Plugin.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if record.Config == nil {
		return map[string]any{}, nil
	}

	return record.Config, nil
}
