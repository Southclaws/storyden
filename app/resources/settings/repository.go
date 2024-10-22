package settings

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/ent"
)

// StorydenPrimarySettingsKey is the key used to store the primary settings data
// in the database. The settings table itself can contain other settings and is
// treated as a key-value store. Storyden itself only cares about the row with
// this key, other rows may be used by plugins or any other integrated systems.
const StorydenPrimarySettingsKey = "storyden_system"

type SettingsRepository struct {
	log *zap.Logger
	db  *ent.Client

	// cached stores the most recent copy of all the settings from the database.
	// Directly changing settings via external database queries will result in
	// settings not immediately updating so it's advised to always go via API.
	cachedSettings *xsync.MapOf[string, any]
}

func New(ctx context.Context, lc fx.Lifecycle, log *zap.Logger, db *ent.Client) (*SettingsRepository, error) {
	d := &SettingsRepository{
		log:            log,
		db:             db,
		cachedSettings: xsync.NewMapOf[string, any](),
	}

	lc.Append(fx.StartHook(func() error {
		if err := d.initDefaults(ctx); err != nil {
			return fault.Wrap(err,
				fctx.With(ctx),
				fmsg.With("failed to initialise default settings"))
		}

		return nil
	}))

	return d, nil
}

// initDefaults is one of the only SettingsRepository writes that happens on first boot. It sets
// up some basic configuration settings for a brand new empty installation.
func (d *SettingsRepository) initDefaults(ctx context.Context) error {
	_, err := d.db.Setting.Get(ctx, StorydenPrimarySettingsKey)
	if ent.IsNotFound(err) {
		_, err = d.setDefaults(ctx)
	}
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (d *SettingsRepository) Get(ctx context.Context) (*Settings, error) {
	r, err := d.db.Setting.Get(ctx, StorydenPrimarySettingsKey)
	if ent.IsNotFound(err) {
		// Ensure defaults are written to the database if they don't exist.
		// This should only happen in tests where initDefaults isn't called.
		return d.setDefaults(ctx)
	}
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapSettings(r)
}

// Set will merge a partial update into the current settings and save new data.
func (d *SettingsRepository) Set(ctx context.Context, s Settings) (*Settings, error) {
	current, err := d.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = current.Merge(s)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	b, err := json.Marshal(current)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := d.db.Setting.
		UpdateOneID(StorydenPrimarySettingsKey).
		SetValue(string(b)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapSettings(r)
}

func (d *SettingsRepository) setDefaults(ctx context.Context) (*Settings, error) {
	b, err := json.Marshal(DefaultSettings)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s, err := d.db.Setting.Create().
		SetID(StorydenPrimarySettingsKey).
		SetValue(string(b)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapSettings(s)
}

// NOTE: There's currently no way to reset/delete or work with non-system data.
