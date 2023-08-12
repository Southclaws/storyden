package reply

import (
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type Reply struct {
	ID post.ID

	Body           string
	Short          string
	Author         profile.Profile
	RootPostID     post.ID
	RootThreadMark string
	ReplyTo        opt.Optional[post.ID]
	Reacts         []*react.React
	Meta           map[string]any
	Assets         []*asset.Asset

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

func (*Reply) GetResourceName() string { return "post" }

func (p Reply) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, post.MakeShortBody(p.Body))
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

	pro, err := profile.FromModel(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	replyTo := replyTo(m)

	return &Reply{
		ID: post.ID(m.ID),

		Body:    m.Body,
		Short:   m.Short,
		Author:  *pro,
		ReplyTo: replyTo,
		Reacts:  dt.Map(m.Edges.Reacts, react.FromModel),
		Meta:    m.Metadata,
		Assets:  dt.Map(m.Edges.Assets, asset.FromModel),

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: opt.NewPtr(m.DeletedAt),
	}, nil
}
