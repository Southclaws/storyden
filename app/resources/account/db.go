package account

import (
	"context"

	"github.com/Southclaws/dt"
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

func (d *database) Create(ctx context.Context, handle string, opts ...option) (*Account, error) {
	withrequired := Account{
		Handle: handle,
		Name:   handle, // default display name is just the handle
	}

	for _, v := range opts {
		v(&withrequired)
	}

	u, err := d.db.Account.
		Create().
		SetHandle(withrequired.Handle).
		SetName(withrequired.Name).
		SetNillableBio(utils.OptionalToPointer(withrequired.Bio)).
		SetNillableID(utils.OptionalID(xid.ID(withrequired.ID))).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(*u), nil
}

func (d *database) GetByID(ctx context.Context, id AccountID) (*Account, error) {
	account, err := d.db.Account.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// threads, posts, err := d.getPostCounts(ctx, account.ID)
	// if err != nil {
	// 	return nil, err
	// }
	acc := FromModel(*account)

	// u.ThreadCount = threads
	// u.PostCount = posts

	return acc, nil
}

func (d *database) GetByHandle(ctx context.Context, handle string) (*Account, error) {
	account, err := d.db.Account.Query().Where(
		account.Handle(handle),
	).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// threads, posts, err := d.getPostCounts(ctx, account.ID)
	// if err != nil {
	// 	return nil, err
	// }
	acc := FromModel(*account)

	// u.ThreadCount = threads
	// u.PostCount = posts

	return acc, nil
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

func (d *database) List(ctx context.Context, sort string, limit, offset int) ([]*Account, error) {
	users, err := d.db.Account.
		Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Asc(account.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return dt.Map(
		dt.Map(users, utils.Deref[ent.Account]),
		utils.ToMap(FromModel),
	), nil
}

func (d *database) Update(ctx context.Context, id AccountID, opts ...Mutation) (*Account, error) {
	update := d.db.Account.UpdateOneID(xid.ID(id))

	for _, fn := range opts {
		fn(update)
	}

	acc, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(*acc), nil
}
