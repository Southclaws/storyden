package reply

import (
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type Reply struct {
	ID post.ID

	Body            string
	Short           string
	Author          profile.Profile
	RootPostID      post.ID
	RootThreadMark  string
	RootThreadTitle string
	ReplyTo         opt.Optional[post.ID]
	Reacts          []*react.React
	Meta            map[string]any
	Assets          []*asset.Asset
	Links           []*datagraph.Link
	URL             opt.Optional[string]

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

func (*Reply) GetResourceName() string { return "post" }

func (r *Reply) GetID() xid.ID           { return xid.ID(r.ID) }
func (r *Reply) GetKind() datagraph.Kind { return datagraph.KindReply }
func (r *Reply) GetName() string         { return r.Short }
func (r *Reply) GetText() string         { return r.Body }
func (r *Reply) GetProps() any           { return r.Meta }

func (p Reply) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, p.Short)
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
		Links:   dt.Map(m.Edges.Links, datagraph.LinkFromModel),

		RootPostID:      post.ID(m.RootPostID),
		RootThreadMark:  opt.NewPtr(m.Edges.Root).OrZero().Slug,
		RootThreadTitle: opt.NewPtr(m.Edges.Root).OrZero().Title,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: opt.NewPtr(m.DeletedAt),
	}, nil
}
