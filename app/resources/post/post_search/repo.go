package post_search

import (
	"context"

	"github.com/Southclaws/dt"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type kindEnum string

const (
	kindThread kindEnum = "thread"
	kindPost   kindEnum = "post"
)

type Filter func(*ent.PostQuery)

type Repository interface {
	Search(ctx context.Context, params pagination.Parameters, filters ...Filter) (*pagination.Result[*post.Post], error)
	GetMany(ctx context.Context, id ...post.ID) ([]*post.Post, error)
}

func WithKinds(ks ...Kind) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			ent_post.Or(
				dt.Map(ks, func(k Kind) predicate.Post {
					switch k {
					case KindThread:
						return ent_post.RootPostIDIsNil()

					case KindPost:
						return ent_post.RootPostIDNotNil()

					default:
						return nil
					}
				})...,
			),
		)
	}
}

func WithKeywords(q string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			ent_post.Or(
				ent_post.And(
					ent_post.RootPostIDIsNil(),
					ent_post.TitleContainsFold(q),
				),
				ent_post.BodyContainsFold(q),
			))
	}
}

func WithAuthorHandle(handle string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			ent_post.HasAuthorWith(
				account.Handle(handle),
			),
		)
	}
}
