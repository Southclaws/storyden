package account_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
)

type Querier struct {
	fx.In

	Ent *ent.Client
}

func (d *Querier) GetByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	q := d.Ent.Account.
		Query().
		Where(account_ent.ID(xid.ID(id))).
		WithTags().
		WithEmails().
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return account.MapAccount(result)
}

func (d *Querier) LookupByHandle(ctx context.Context, handle string) (*account.Account, bool, error) {
	q := d.Ent.Account.
		Query().
		Where(account_ent.Handle(handle)).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// threads, posts, err := d.getPostCounts(ctx, account.ID)
	// if err != nil {
	// 	return nil, err
	// }
	acc, err := account.MapAccount(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	// u.ThreadCount = threads
	// u.PostCount = posts

	return acc, true, nil
}
