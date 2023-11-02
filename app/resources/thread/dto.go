package thread

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/utils"
)

type Thread struct {
	ID        post.ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]

	Title       string
	Slug        string
	Short       string
	Pinned      bool
	Author      profile.Profile
	Tags        []string
	Category    category.Category
	Status      post.Status
	Posts       []*reply.Reply
	Reacts      []*react.React
	Meta        map[string]any
	Assets      []*asset.Asset
	Collections []*collection.Collection
	Links       link.Links
}

func (*Thread) GetResourceName() string { return "thread" }

func FromModel(m *ent.Post) (*Thread, error) {
	authorEdge, err := m.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.FromModel(authorEdge)
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

	// Thread data structure will always contain one post: itself in post form.
	first, err := transform(m)
	if err != nil {
		return nil, err
	}
	posts := []*reply.Reply{first}

	if p, err := m.Edges.PostsOrErr(); err == nil && len(p) > 0 {
		transformed, err := dt.MapErr(p, transform)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		posts = append(posts, transformed...)
	}

	collectionsEdge := opt.NewIf(m.Edges.Collections, func(c []*ent.Collection) bool { return c != nil })

	collections, err := opt.MapErr(collectionsEdge, func(c []*ent.Collection) ([]*collection.Collection, error) {
		out, err := dt.MapErr(c, collection.FromModel)
		if err != nil {
			return nil, fault.Wrap(err)
		}
		return out, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Thread{
		ID:        post.ID(m.ID),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: opt.NewPtr(m.DeletedAt),

		Title:       m.Title,
		Slug:        m.Slug,
		Short:       m.Short,
		Pinned:      m.Pinned,
		Author:      *pro,
		Tags:        dt.Map(m.Edges.Tags, func(t *ent.Tag) string { return t.Name }),
		Category:    utils.Deref(category.FromModel(m.Edges.Category)),
		Status:      post.NewStatusFromEnt(m.Status),
		Posts:       posts,
		Reacts:      dt.Map(m.Edges.Reacts, react.FromModel),
		Meta:        m.Metadata,
		Assets:      dt.Map(m.Edges.Assets, asset.FromModel),
		Collections: collections.OrZero(),
		Links:       dt.Map(m.Edges.Links, link.Map),
	}, nil
}
