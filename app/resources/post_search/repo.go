package post_search

import (
	"context"

	"github.com/Southclaws/dt"

	post_resource "github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/post"
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
	Search(ctx context.Context, opts ...Filter) ([]*post_resource.Reply, error)
}

func WithKinds(ks ...Kind) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			post.Or(
				dt.Map(ks, func(k Kind) predicate.Post {
					switch k {
					case KindThread:
						return post.First(true)

					case KindPost:
						return post.First(false)

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
		pq.Where(post.And(
			post.First(true),
			post.TitleContainsFold(q),
		))
	}
}

func WithBodyContains(q string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(post.BodyContains(q))
	}
}

func WithAuthorHandle(handle string) Filter {
	return func(pq *ent.PostQuery) {
		pq.Where(
			post.HasAuthorWith(
				account.Handle(handle),
			),
		)
	}
}
