package invitation_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/invitation"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (d *Writer) Create(ctx context.Context, creator account.AccountID, message opt.Optional[string]) (*invitation.Invitation, error) {
	create := d.db.Invitation.Create()

	create.SetCreatorAccountID(xid.ID(creator))
	create.SetNillableMessage(message.Ptr())

	result, err := create.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	result.Edges.Creator, err = result.QueryCreator().Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	inv, err := invitation.Map(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return inv, nil
}

func (d *Writer) Delete(ctx context.Context, id xid.ID) error {
	err := d.db.Invitation.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
