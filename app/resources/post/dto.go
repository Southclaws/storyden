package post

import (
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type PostID xid.ID

func (u PostID) String() string { return xid.ID(u).String() }

type Post struct {
	ID PostID

	Body           string
	Short          string
	Author         Author
	RootPostID     PostID
	RootThreadMark string
	ReplyTo        opt.Optional[PostID]
	Reacts         []*react.React
	Meta           map[string]any
	Assets         []*asset.Asset

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt opt.Optional[time.Time]
}

func (*Post) GetResourceName() string { return "post" }

type Author struct {
	ID        account.AccountID
	Name      string
	Handle    string
	Admin     bool
	CreatedAt time.Time
}

func (p Post) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, MakeShortBody(p.Body))
}

func replyTo(m *ent.Post) opt.Optional[PostID] {
	if m.Edges.ReplyTo != nil {
		return opt.New(PostID(m.Edges.ReplyTo.ID))
	}

	return opt.NewEmpty[PostID]()
}

func FromModel(m *ent.Post) (w *Post) {
	replyTo := replyTo(m)

	return &Post{
		ID: PostID(m.ID),

		Body:  m.Body,
		Short: m.Short,
		Author: Author{
			ID:        account.AccountID(m.Edges.Author.ID),
			Name:      m.Edges.Author.Name,
			Handle:    m.Edges.Author.Handle,
			Admin:     m.Edges.Author.Admin,
			CreatedAt: m.Edges.Author.CreatedAt,
		},
		ReplyTo: replyTo,
		Reacts:  dt.Map(m.Edges.Reacts, react.FromModel),
		Meta:    m.Metadata,
		Assets:  dt.Map(m.Edges.Assets, asset.FromModel),

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: opt.NewPtr(m.DeletedAt),
	}
}
