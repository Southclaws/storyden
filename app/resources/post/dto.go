package post

import (
	"fmt"
	"time"

	"4d63.com/optional"
	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type PostID xid.ID

func (u PostID) String() string { return xid.ID(u).String() }

type Post struct {
	ID PostID `json:"id"`

	Body           string
	Short          string
	Author         Author
	RootPostID     PostID
	RootThreadMark string
	ReplyTo        optional.Optional[PostID]
	Reacts         []*react.React
	Meta           map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt optional.Optional[time.Time]
}

func (*Post) GetResourceName() string { return "post" }

type Author struct {
	ID        account.AccountID `json:"id"`
	Name      string            `json:"name"`
	Handle    string
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p Post) String() string {
	return fmt.Sprintf("post %s by '%s' at %s\n'%s'", p.ID.String(), p.Author.Handle, p.CreatedAt, MakeShortBody(p.Body))
}

func replyTo(m *ent.Post) optional.Optional[PostID] {
	if m.Edges.ReplyTo != nil {
		return optional.Of(PostID(m.Edges.ReplyTo.ID))
	}

	return optional.Empty[PostID]()
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

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: optional.OfPtr(m.DeletedAt),
	}
}
