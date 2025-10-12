package notify_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (n *Writer) Notification(ctx context.Context,
	accountID account.AccountID,
	event notification.Event,
	item opt.Optional[datagraph.ItemRef],
	source opt.Optional[account.AccountID],
) (*notification.NotificationRef, error) {
	create := n.db.Notification.Create()

	create.SetOwnerID(xid.ID(accountID)).
		SetEventType(event.String()).
		SetRead(false)

	if i, ok := item.Get(); ok {
		create.
			SetDatagraphKind(i.GetKind().String()).
			SetDatagraphID(i.GetID())
	}

	source.Call(func(value account.AccountID) { create.SetSourceAccountID(xid.ID(value)) })

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nr, err := notification.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nr, nil
}

func (n *Writer) SetRead(ctx context.Context, id xid.ID, read bool) (*notification.NotificationRef, error) {
	r, err := n.db.Notification.UpdateOneID(id).
		SetRead(read).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nr, err := notification.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nr, nil
}
