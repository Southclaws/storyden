package tag_writer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db}
}

func (w *Writer) Add(ctx context.Context, names ...tag_ref.Name) ([]*tag_ref.Tag, error) {
	nameStrings := tag_ref.Names(names).Strings()

	newTags := dt.Map(nameStrings, func(n string) *ent.TagCreate {
		return w.db.Tag.Create().SetName(n)
	})

	create := w.db.Tag.
		CreateBulk(newTags...).
		OnConflictColumns(ent_tag.FieldName).
		DoNothing()

	err := create.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := w.db.Tag.Query().Where(ent_tag.NameIn(nameStrings...)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := dt.Map(r, tag_ref.Map(nil))

	return tags, nil
}

func (w *Writer) Remove(ctx context.Context, names ...tag_ref.Name) error {
	nameStrings := tag_ref.Names(names).Strings()

	_, err := w.db.Tag.Delete().
		Where(ent_tag.NameIn(nameStrings...)).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
