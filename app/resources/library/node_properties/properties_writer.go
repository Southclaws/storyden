package node_properties

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/property"
	"github.com/rs/xid"
)

type Writer struct {
	db *ent.Client
}

func (w *Writer) Update(ctx context.Context, nid library.NodeID, schema library.PropertySchema, props library.ExistingPropertyMutations) (*library.PropertyTable, error) {
	tx, err := w.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		err = tx.Rollback()
	}()

	updated := library.PropertyTable{
		Schema: schema,
	}

	for _, prop := range props {
		create := tx.Property.Create().
			SetValue(prop.Value).
			SetFieldID(prop.ID).
			SetNodeID(xid.ID(nid))

		create.OnConflictColumns(property.FieldFieldID, property.FieldNodeID).UpdateNewValues()

		r, err := create.Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		updated.Properties = append(updated.Properties, &library.Property{
			Field: prop.PropertySchemaField,
			Value: opt.New(r.Value),
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &updated, nil
}
