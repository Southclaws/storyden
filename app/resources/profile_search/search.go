package profile_search

import (
	"context"
	"math"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
)

type Filter func(*ent.AccountQuery)

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Profiles    []*profile.Profile
}

type Repository interface {
	Search(ctx context.Context, page int, pageSize int, opts ...Filter) (*Result, error)
}

func WithDisplayNameContains(q string) Filter {
	return func(pq *ent.AccountQuery) {
		pq.Where(account.And(
			account.NameContainsFold(q),
		))
	}
}

func WithHandleContains(q string) Filter {
	return func(pq *ent.AccountQuery) {
		pq.Where(account.And(
			account.HandleContainsFold(q),
		))
	}
}

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Search(ctx context.Context, page int, size int, filters ...Filter) (*Result, error) {
	total, err := d.db.Account.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	q := d.db.Account.Query().
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Asc(account.FieldCreatedAt))

	for _, fn := range filters {
		fn(q)
	}

	r, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nextPage := opt.NewSafe(page+1, len(r) >= size)

	if len(r) > 1 {
		r = r[:len(r)-1]
	}

	profiles, err := dt.MapErr(r, profile.FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Result{
		PageSize:    size,
		Results:     len(profiles),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Profiles:    profiles,
	}, nil
}
