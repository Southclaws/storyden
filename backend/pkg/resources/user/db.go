package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/user"
	"github.com/Southclaws/storyden/backend/internal/utils"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) CreateUser(ctx context.Context, email string, username string) (*User, error) {
	u, err := d.db.User.
		Create().
		SetEmail(email).
		SetName(username).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Ref(FromModel(*u)), nil
}

func (d *database) GetUser(ctx context.Context, userId UserID, public bool) (*User, error) {
	user, err := d.db.User.Get(ctx, uuid.UUID(userId))
	if err != nil {
		if model.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// threads, posts, err := d.getUserPostCounts(ctx, user.ID)
	// if err != nil {
	// 	return nil, err
	// }
	u := FromModel(*user)

	// u.ThreadCount = threads
	// u.PostCount = posts

	return &u, nil
}

// func (d *database) getUserPostCounts(ctx context.Context, id string) (int, int, error) {
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
// 			from "User" u
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

func (d *database) GetUserByEmail(ctx context.Context, email string, public bool) (*User, error) {
	user, err := d.db.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return utils.Ref(FromModel(*user)), nil
}

func (d *database) GetUsers(ctx context.Context, sort string, limit, offset int, public bool) ([]User, error) {
	users, err := d.db.User.
		Query().
		Limit(limit).
		Offset(offset).
		Order(model.Asc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if public {
		return lo.Map(
			lo.Map(users, utils.Deref[model.User]),
			utils.ToMap(FromModelPublic),
		), nil
	}

	return lo.Map(
		lo.Map(users, utils.Deref[model.User]),
		utils.ToMap(FromModel),
	), nil
}

func (d *database) UpdateUser(ctx context.Context, userId UserID, email, name, bio *string) (*User, error) {
	update := d.db.User.UpdateOneID(uuid.UUID(userId))

	// TODO: This is awful. Make this more ergonomic.
	if email != nil {
		update.SetEmail(*email)
	}
	if name != nil {
		update.SetEmail(*name)
	}
	if bio != nil {
		update.SetEmail(*bio)
	}

	u, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Ref(FromModel(*u)), nil
}

func (d *database) SetAdmin(ctx context.Context, userId UserID, status bool) error {
	_, err := d.db.User.
		UpdateOneID(uuid.UUID(userId)).
		SetAdmin(status).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *database) Ban(ctx context.Context, userId UserID) (*User, error) {
	u, err := d.db.User.
		UpdateOneID(uuid.UUID(userId)).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Ref(FromModel(*u)), nil
}

func (d *database) Unban(ctx context.Context, userId UserID) (*User, error) {
	u, err := d.db.User.
		UpdateOneID(uuid.UUID(userId)).
		ClearDeletedAt().
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Ref(FromModel(*u)), nil
}
