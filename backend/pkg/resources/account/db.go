package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/backend/internal/utils"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context, email string, username string) (*Account, error) {
	u, err := d.db.Account.
		Create().
		SetEmail(email).
		SetName(username).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Ref(FromModel(*u)), nil
}

func (d *database) GetByID(ctx context.Context, userId AccountID) (*Account, error) {
	account, err := d.db.Account.Get(ctx, uuid.UUID(userId))
	if err != nil {
		if model.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
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

func (d *database) LookupByEmail(ctx context.Context, email string) (*Account, bool, error) {
	account, err := d.db.Account.
		Query().
		Where(account.Email(email)).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return utils.Ref(FromModel(*account)), true, nil
}

func (d *database) List(ctx context.Context, sort string, limit, offset int) ([]Account, error) {
	users, err := d.db.Account.
		Query().
		Limit(limit).
		Offset(offset).
		Order(model.Asc(account.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(
		lo.Map(users, utils.Deref[model.Account]),
		utils.ToMap(FromModel),
	), nil
}
