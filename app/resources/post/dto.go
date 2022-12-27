package post

import (
	"time"

	"4d63.com/optional"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
)

type PostID xid.ID

func (u PostID) String() string { return xid.ID(u).String() }

type Post struct {
	ID PostID `json:"id"`

	Body       string                    `json:"body"`
	Short      string                    `json:"short"`
	Author     Author                    `json:"author"`
	RootPostID PostID                    `json:"rootPostId"`
	ReplyTo    optional.Optional[PostID] `json:"replyTo"`
	Reacts     []react.React             `json:"reacts"`
	Meta       map[string]any

	CreatedAt time.Time                    `json:"createdAt"`
	UpdatedAt time.Time                    `json:"updatedAt"`
	DeletedAt optional.Optional[time.Time] `json:"deletedAt"`
}

const Role = "Post"

func (u *Post) GetRole() string { return Role }

type Author struct {
	ID        account.AccountID `json:"id"`
	Name      string            `json:"name"`
	Handle    string
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"createdAt"`
}

func replyTo(m *ent.Post) optional.Optional[PostID] {
	if m.Edges.ReplyTo != nil {
		return optional.Of(PostID(m.Edges.ReplyTo.ID))
	}

	return optional.Empty[PostID]()
}

func FromModel(m *ent.Post) (w *Post) {
	replyTo := replyTo(m)

	reacts := lo.Map(m.Edges.Reacts, func(t *ent.React, i int) react.React {
		r := react.FromModel(t)
		return *r
	})

	// replyTo := utils.OptionalSlice[xid.ID](m.ReplyToPostID)

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
		RootPostID: PostID(m.RootPostID),
		ReplyTo:    replyTo,
		Reacts:     reacts,
		Meta:       m.Metadata,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: optional.OfPtr(m.DeletedAt),
	}
}
