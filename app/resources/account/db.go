package account

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/utils"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context, handle string, opts ...Option) (*Account, error) {
	withrequired := Account{
		Handle: handle,
		Name:   handle, // default display name is just the handle
	}

	for _, v := range opts {
		v(&withrequired)
	}

	a, err := d.db.Account.
		Create().
		SetHandle(withrequired.Handle).
		SetName(withrequired.Name).
		SetNillableBio(utils.OptionalToPointer(withrequired.Bio)).
		SetNillableID(utils.OptionalID(xid.ID(withrequired.ID))).
		SetAdmin(withrequired.Admin).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(a)
}

func (d *database) GetByID(ctx context.Context, id AccountID) (*Account, error) {
	q := d.db.Account.
		Query().
		Where(account.ID(xid.ID(id))).
		WithTags().
		WithAuthentication()

	account, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(account)
}

func (d *database) LookupByHandle(ctx context.Context, handle string) (*Account, bool, error) {
	q := d.db.Account.
		Query().
		Where(account.Handle(handle)).
		WithAuthentication()

	account, err := q.Only(ctx)
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
	acc, err := FromModel(account)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	// u.ThreadCount = threads
	// u.PostCount = posts

	return acc, true, nil
}

// func (d *database) getPostCounts(ctx context.Context, id string) (int, int, error) {
// 	type R struct {
// 		Threads int `json:"threads"`
// 		Posts   int `json:"posts"`
// 	}
// 	var count []R
// 	err := d.db.Prisma.
// 		QueryRaw(`
// 		select
// 			count(*) filter (where "first") as threads,
// 			count(*) filter (where not "first") as posts
// 		from (
// 			select p.first
// 			from "Account" u
// 			inner join "Post" p on p."userId" = u.id
// 			where u.id = $1
// 		) t`, id).
// 		Exec(ctx, &count)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	if len(count) == 0 {
// 		return 0, 0, nil
// 	}

// 	return count[0].Threads, count[0].Posts, nil
// }

func (d *database) Update(ctx context.Context, id AccountID, opts ...Mutation) (*Account, error) {
	update := d.db.Account.UpdateOneID(xid.ID(id))

	for _, fn := range opts {
		fn(update)
	}

	acc, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(acc)
}
