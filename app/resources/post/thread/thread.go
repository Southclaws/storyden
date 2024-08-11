package thread

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

type Thread struct {
	post.Post

	Title  string
	Slug   string
	Short  string
	Pinned bool

	Replies     []*reply.Reply
	Category    category.Category
	Visibility  visibility.Visibility
	Tags        []string
	Collections []*collection.Collection
	Related     datagraph.ItemList
}

func (*Thread) GetResourceName() string { return "thread" }

func (t *Thread) GetName() string { return t.Title }
func (t *Thread) GetSlug() string { return t.Slug }
func (t *Thread) GetDesc() string { return t.Short }

func FromModel(m *ent.Post) (*Thread, error) {
	categoryEdge, err := m.Edges.CategoryOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	authorEdge, err := m.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	category := category.FromModel(categoryEdge)

	pro, err := profile.ProfileFromModel(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	transform := func(v *ent.Post) (*reply.Reply, error) {
		// hydrate the thread-specific info here. post.FromModel cannot do this
		// as this info is only available in the context of a thread of posts.
		dto, err := reply.FromModel(v)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		dto.RootThreadMark = m.Slug
		dto.RootPostID = post.ID(m.ID)
		return dto, nil
	}

	first, err := transform(m)
	if err != nil {
		return nil, err
	}
	posts := []*reply.Reply{}

	if p, err := m.Edges.PostsOrErr(); err == nil && len(p) > 0 {
		transformed, err := dt.MapErr(p, transform)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		posts = append(posts, transformed...)
	}

	collectionsEdge := opt.NewIf(m.Edges.Collections, func(c []*ent.Collection) bool { return c != nil })

	collections, err := opt.MapErr(collectionsEdge, func(c []*ent.Collection) ([]*collection.Collection, error) {
		out, err := dt.MapErr(c, collection.MapCollection)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		return out, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Thread{
		Post: post.Post{
			ID: post.ID(m.ID),

			Content: first.Post.Content,
			Author:  *pro,
			Reacts:  dt.Map(m.Edges.Reacts, react.FromModel),
			Assets:  dt.Map(m.Edges.Assets, asset.FromModel),
			Links:   dt.Map(m.Edges.Links, datagraph.LinkFromModel),
			Meta:    m.Metadata,

			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: opt.NewPtr(m.DeletedAt),
		},

		Title:  m.Title,
		Slug:   m.Slug,
		Short:  m.Short,
		Pinned: m.Pinned,

		Replies:     posts,
		Category:    *category,
		Visibility:  visibility.NewVisibilityFromEnt(m.Visibility),
		Tags:        dt.Map(m.Edges.Tags, func(t *ent.Tag) string { return t.Name }),
		Collections: collections.OrZero(),
	}, nil
}
