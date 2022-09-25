package account

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/fault/errtag"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/internal/utils"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
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
		if model.IsConstraintError(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.AlreadyExists{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return utils.Ref(FromModel(*u)), nil
}

func (d *database) GetByID(ctx context.Context, userId AccountID) (*Account, error) {
	account, err := d.db.Account.Get(ctx, xid.ID(userId))
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	// threads, posts, err := d.getPostCounts(ctx, account.ID)
	// if err != nil {
	// 	return nil, err
	// }
	u := FromModel(*account)

	// u.ThreadCount = threads
	// u.PostCount = posts

	return &u, nil
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

func (d *database) List(ctx context.Context, sort string, limit, offset int) ([]Account, error) {
	users, err := d.db.Account.
		Query().
		Limit(limit).
		Offset(offset).
		Order(model.Asc(account.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return dt.Map(
		dt.Map(users, utils.Deref[model.Account]),
		utils.ToMap(FromModel),
	), nil
}
