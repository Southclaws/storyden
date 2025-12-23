package post_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
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
	Locate(ctx context.Context, externalID post.ID) (*Location, error)
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
				ent_account.Handle(handle),
			),
		)
	}
}

func WithAuthors(ids ...account.AccountID) Filter {
	return func(pq *ent.PostQuery) {
		if len(ids) == 0 {
			return
		}
		pq.Where(ent_post.AccountPostsIn(dt.Map(ids, func(id account.AccountID) xid.ID {
			return xid.ID(id)
		})...))
	}
}

func WithCategories(ids ...category.CategoryID) Filter {
	return func(pq *ent.PostQuery) {
		if len(ids) == 0 {
			return
		}
		pq.Where(ent_post.CategoryIDIn(dt.Map(ids, func(id category.CategoryID) xid.ID {
			return xid.ID(id)
		})...))
	}
}

func WithTags(names ...tag_ref.Name) Filter {
	return func(pq *ent.PostQuery) {
		if len(names) == 0 {
			return
		}
		predicates := make([]predicate.Post, len(names))
		for i, name := range names {
			predicates[i] = ent_post.HasTagsWith(ent_tag.NameEQ(name.String()))
		}
		pq.Where(ent_post.And(predicates...))
	}
}
