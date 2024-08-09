package post_search

import (
	"context"

	"github.com/Southclaws/dt"

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
	Search(ctx context.Context, opts ...Filter) ([]*post.Post, error)
	GetMany(ctx context.Context, id ...post.ID) ([]*post.Post, error)
}

func WithKinds(ks ...Kind) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			ent_post.Or(
				dt.Map(ks, func(k Kind) predicate.Post {
					switch k {
					case KindThread:
						return ent_post.First(true)

					case KindPost:
						return ent_post.First(false)

					default:
						return nil
					}
				})...,
			),
		)
	}
}

func WithTitleContains(q string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(ent_post.And(
			ent_post.First(true),
			ent_post.TitleContainsFold(q),
		))
	}
}

func WithBodyContains(q string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(ent_post.BodyContains(q))
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
