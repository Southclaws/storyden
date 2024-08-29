package reply

import (
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type Reply struct {
	post.Post

	RootPostID      post.ID
	RootThreadMark  string
	RootThreadTitle string
	ReplyTo         opt.Optional[post.ID]
}

func (*Reply) GetResourceName() string { return "post" }

func (r *Reply) GetName() string {
	if xid.ID(r.RootPostID).IsZero() {
		return r.RootThreadTitle
	}

	return fmt.Sprintf("reply to: %s", r.RootThreadTitle)
}

func (r *Reply) GetSlug() string {
	if xid.ID(r.RootPostID).IsZero() {
		return r.ID.String()
	}
	return r.RootThreadMark
}
func (r *Reply) GetDesc() string { return r.Content.Short() }

func (p Reply) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, p.Content.Short())
}

func replyTo(m *ent.Post) opt.Optional[post.ID] {
	if m.Edges.ReplyTo != nil {
		return opt.New(post.ID(m.Edges.ReplyTo.ID))
	}

	return opt.NewEmpty[post.ID]()
}

func FromModel(m *ent.Post) (*Reply, error) {
	authorEdge, err := m.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.ProfileFromModel(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	content, err := content.NewRichText(m.Body)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	replyTo := replyTo(m)

	var rootPostID post.ID
	var rootThreadMark string
	var rootThreadTitle string
	if m.RootPostID == xid.NilID() {
		// A root post was passed, which is still valid in some cases.
		rootThreadMark = m.Slug
		rootThreadTitle = m.Title
	} else {
		rootPostID = post.ID(m.RootPostID)
		rootThreadMark = opt.NewPtr(m.Edges.Root).OrZero().Slug
		rootThreadTitle = opt.NewPtr(m.Edges.Root).OrZero().Title
	}

	link := opt.Map(opt.NewPtr(m.Edges.Link), func(in ent.Link) link_ref.LinkRef {
		return *link_ref.Map(&in)
	})

	return &Reply{
		Post: post.Post{
			ID: post.ID(m.ID),

			Content: content,
			Author:  *pro,
			Reacts:  dt.Map(m.Edges.Reacts, react.FromModel),
			Assets:  dt.Map(m.Edges.Assets, asset.FromModel),
			WebLink: link,
			Meta:    m.Metadata,

			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: opt.NewPtr(m.DeletedAt),
		},
		ReplyTo: replyTo,

		RootPostID:      rootPostID,
		RootThreadMark:  rootThreadMark,
		RootThreadTitle: rootThreadTitle,
	}, nil
}
