package invitation_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/invitation"
	"github.com/Southclaws/storyden/internal/ent"
	invitation_ent "github.com/Southclaws/storyden/internal/ent/invitation"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

type Filter func(*ent.InvitationQuery)

func WithCreator(id account.AccountID) Filter {
	return func(q *ent.InvitationQuery) {
		q.Where(invitation_ent.CreatorAccountIDEQ(xid.ID(id)))
	}
}

func (d *Querier) GetByID(ctx context.Context, id xid.ID) (*invitation.Invitation, error) {
	q := d.db.Invitation.
		Query().
		Where(invitation_ent.ID(id)).
		WithCreator()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	inv, err := invitation.Map(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return inv, nil
}

func (d *Querier) List(ctx context.Context, opts ...Filter) ([]*invitation.Invitation, error) {
	q := d.db.Invitation.Query()

	for _, opt := range opts {
		opt(q)
	}

	result, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	invs, err := dt.MapErr(result, invitation.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return invs, nil
}
