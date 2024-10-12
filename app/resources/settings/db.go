package settings

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/setting"
)

type database struct {
	db *ent.Client
}

func New(ctx context.Context, lc fx.Lifecycle, db *ent.Client) (Repository, error) {
	d := &database{db}

	lc.Append(fx.StartHook(func() error {
		if err := d.Init(ctx); err != nil {
			return fault.Wrap(err)
		}

		return nil
	}))

	return d, nil
}

// Init is one of the only database writes that happens on first boot. It sets
// up some basic configuration settings for a brand new empty installation.
func (d *database) Init(ctx context.Context) error {
	r, err := d.db.Setting.Query().Count(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if r > 0 {
		return nil
	}

	if err := d.SetValue(ctx, "Title", DefaultTitle); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := d.SetValue(ctx, "Description", DefaultDescription); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := d.SetValue(ctx, "Content", DefaultContent); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := d.SetValue(ctx, "AccentColour", DefaultColour); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (d *database) Get(ctx context.Context) (*Settings, error) {
	r, err := d.db.Setting.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return fromEnt(r)
}

func (d *database) Set(ctx context.Context, s Partial) (*Settings, error) {
	kvs, err := toEnt(s)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	for _, entry := range kvs {
		if err := d.SetValue(ctx, entry.key, entry.value); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return d.Get(ctx)
}

func (d *database) SetValue(ctx context.Context, key, value string) error {
	u := d.db.Setting.
		Create().
		SetID(key).
		SetValue(value).
		OnConflict(
			sql.ConflictColumns(setting.FieldID),
			sql.ResolveWithNewValues(),
		).
		SetValue(value)
	if err := u.Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (d *database) GetValue(ctx context.Context, key string) (string, error) {
	s, err := d.db.Setting.Get(ctx, key)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return s.Value, nil
}
