package thread

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/utils"
)

type Thread struct {
	ID        post.PostID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]

	Title    string
	Slug     string
	Short    string
	Pinned   bool
	Author   post.Author
	Tags     []string
	Category category.Category
	Status   Status
	Posts    []*post.Post
	Reacts   []*react.React
	Meta     map[string]any
}

func (*Thread) GetResourceName() string { return "thread" }

func FromModel(m *ent.Post) *Thread {
	transform := func(v *ent.Post) *post.Post {
		// hydrate the thread-specific info here. post.FromModel cannot do this
		// as this info is only available in the context of a thread of posts.
		dto := post.FromModel(v)
		dto.RootThreadMark = m.Slug
		dto.RootPostID = post.PostID(m.ID)
		return dto
	}

	// Thread data structure will always contain one post: itself in post form.
	posts := []*post.Post{
		transform(m),
	}

	if p, err := m.Edges.PostsOrErr(); err == nil && len(p) > 0 {
		posts = append(posts, dt.Map(p, transform)...)
	}

	return &Thread{
		ID:        post.PostID(m.ID),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: opt.NewPtr(m.DeletedAt),

		Title:  m.Title,
		Slug:   m.Slug,
		Short:  m.Short,
		Pinned: m.Pinned,
		Author: post.Author{
			ID:     account.AccountID(m.Edges.Author.ID),
			Handle: m.Edges.Author.Handle,
			Name:   m.Edges.Author.Name,
		},
		Tags:     dt.Map(m.Edges.Tags, func(t *ent.Tag) string { return t.Name }),
		Category: utils.Deref(category.FromModel(m.Edges.Category)),
		Status:   NewStatusFromEnt(m.Status),
		Posts:    posts,
		Reacts:   dt.Map(m.Edges.Reacts, react.FromModel),
		Meta:     m.Metadata,
	}
}
