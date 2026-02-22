package notify_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entnotification "github.com/Southclaws/storyden/internal/ent/notification"
)

type Writer struct {
	db          *ent.Client
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, roleQuerier *role_querier.Querier) *Writer {
	return &Writer{db: db, roleQuerier: roleQuerier}
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

	r, err = n.db.Notification.Query().
		Where(entnotification.ID(r.ID)).
		WithSource().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if source := r.Edges.Source; source != nil {
		if err := n.roleQuerier.HydrateRoleEdges(ctx, source); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
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

	r, err = n.db.Notification.Query().
		Where(entnotification.ID(r.ID)).
		WithSource().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if source := r.Edges.Source; source != nil {
		if err := n.roleQuerier.HydrateRoleEdges(ctx, source); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	nr, err := notification.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nr, nil
}

func (n *Writer) UpdateStatusMany(ctx context.Context, accountID account.AccountID, notifications []*notification.NotificationRef) ([]*notification.NotificationRef, error) {
	tx, err := n.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		_ = tx.Rollback()
	}()

	updatedIDs := make([]xid.ID, 0, len(notifications))

	for _, notif := range notifications {
		r, err := tx.Notification.UpdateOneID(xid.ID(notif.ID)).
			Where(entnotification.HasOwnerWith(entaccount.ID(xid.ID(accountID)))).
			SetRead(notif.Read).
			Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		updatedIDs = append(updatedIDs, r.ID)
	}
	if len(updatedIDs) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return []*notification.NotificationRef{}, nil
	}

	rows, err := tx.Notification.Query().
		Where(entnotification.IDIn(updatedIDs...)).
		WithSource().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sources := make([]*ent.Account, 0, len(rows))
	for _, row := range rows {
		if source := row.Edges.Source; source != nil {
			sources = append(sources, source)
		}
	}
	if err := n.roleQuerier.HydrateRoleEdges(ctx, sources...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rowByID := make(map[xid.ID]*ent.Notification, len(rows))
	for _, row := range rows {
		rowByID[row.ID] = row
	}

	updated := make([]*notification.NotificationRef, 0, len(updatedIDs))
	for _, updatedID := range updatedIDs {
		row := rowByID[updatedID]
		if row == nil {
			return nil, fault.New("updated notification missing after requery", fctx.With(ctx))
		}

		nr, err := notification.Map(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		updated = append(updated, nr)
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return updated, nil
}
