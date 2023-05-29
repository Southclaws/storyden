package settings

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/setting"
)

type database struct {
	db *ent.Client
}

func New(ctx context.Context, db *ent.Client) (Repository, error) {
	d := &database{db}

	if err := d.Init(ctx); err != nil {
		return nil, fault.Wrap(err)
	}

	return d, nil
}

func (d *database) Init(ctx context.Context) error {
	r, err := d.db.Setting.Query().Count(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if r > 0 {
		return nil
	}

	if err := d.SetValue(ctx, "Title", "Storyden"); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := d.SetValue(ctx, "Description", "A forum for the modern age"); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := d.SetValue(ctx, "AccentColour", "hsl(157, 65%, 44%)"); err != nil {
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
