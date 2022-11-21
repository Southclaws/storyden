package react

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type ReactID xid.ID

type React struct {
	ID     ReactID `json:"id"`
	Emoji  string  `json:"emoji"`
	UserID string  `json:"user"`
	PostID string  `json:"post"`
}

func FromModel(ent *ent.React) *React {
	return &React{
		ID:     ReactID(ent.ID),
		Emoji:  ent.Emoji,
		UserID: ent.Edges.Account.ID.String(),
		PostID: ent.Edges.Post.ID.String(),
	}
}

func IsValidEmoji(e string) (string, bool) {
	if len(e) == 0 {
		return "", false
	}

	if e[1] == ':' {
		return "", false
	}

	return e, gomoji.ContainsEmoji(e)
}
