package moderation_note

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type Note struct {
	ID        ID
	Author    opt.Optional[profile.Ref]
	Content   string
	CreatedAt time.Time
}

type Notes []*Note

func Map(note *ent.ModerationNote) (*Note, error) {
	authorEdge := opt.NewPtr(note.Edges.Author)

	author, err := opt.MapErr(authorEdge, func(a ent.Account) (profile.Ref, error) {
		p, err := profile.MapRef(&a)
		if err != nil {
			return profile.Ref{}, err
		}
		return *p, nil
	})
	if err != nil {
		return nil, err
	}

	return &Note{
		ID:        note.ID,
		Author:    author,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
	}, nil
}
