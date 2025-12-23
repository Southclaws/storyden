package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type postSearcher struct {
	post_search post_search.Repository
}

func (s *postSearcher) Search(ctx context.Context, query string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	o := []post_search.Filter{
		post_search.WithKeywords(query),
	}

	opts.Kinds.Call(func(value []datagraph.Kind) {
		ks := []post_search.Kind{}
		for _, k := range value {
			switch k {
			case datagraph.KindThread:
				ks = append(ks, post_search.KindThread)
			case datagraph.KindReply:
				ks = append(ks, post_search.KindPost)
			}
		}

		o = append(o, post_search.WithKinds(ks...))
	})

	opts.Authors.Call(func(value []account.AccountID) {
		o = append(o, post_search.WithAuthors(value...))
	})

	opts.Categories.Call(func(value []category.CategoryID) {
		o = append(o, post_search.WithCategories(value...))
	})

	opts.Tags.Call(func(value []tag_ref.Name) {
		o = append(o, post_search.WithTags(value...))
	})

	rs, err := s.post_search.Search(ctx, p, o...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := dt.MapErr(rs.Items, func(r *post.Post) (datagraph.Item, error) {
		if r.ID == r.Root {
			return s.mapToThread(r)
		}
		return s.mapToReply(r)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.ConvertPageResult(*rs, items)

	return &result, nil
}

func (s *postSearcher) mapToThread(p *post.Post) (datagraph.Item, error) {
	content, err := datagraph.NewRichText(p.Content.HTML())
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &thread.Thread{
		Post: post.Post{
			ID:        p.ID,
			Root:      p.Root,
			Content:   content,
			Author:    p.Author,
			Likes:     p.Likes,
			Reacts:    p.Reacts,
			Assets:    p.Assets,
			WebLink:   p.WebLink,
			Meta:      p.Meta,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt,
		},
		Title: p.Title,
		Slug:  p.Slug,
		Short: p.Content.Short(),
	}, nil
}

func (s *postSearcher) mapToReply(p *post.Post) (datagraph.Item, error) {
	content, err := datagraph.NewRichText(p.Content.HTML())
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &reply.Reply{
		Post: post.Post{
			ID:        p.ID,
			Root:      p.Root,
			Content:   content,
			Author:    p.Author,
			Likes:     p.Likes,
			Reacts:    p.Reacts,
			Assets:    p.Assets,
			WebLink:   p.WebLink,
			Meta:      p.Meta,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt,
		},
		RootPostID: p.Root,
		Slug:       p.Slug,
	}, nil
}

func (s *postSearcher) MatchFast(ctx context.Context, q string, limit int, opts searcher.Options) (datagraph.MatchList, error) {
	return nil, searcher.ErrFastMatchesUnavailable
}
