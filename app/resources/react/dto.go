package react

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
)

type ReactID xid.ID

type React struct {
	ID     ReactID `json:"id"`
	Emoji  string  `json:"emoji"`
	UserID string  `json:"user"`
	PostID string  `json:"post"`
}

func FromModel(model *model.React) *React {
	return &React{
		ID:     ReactID(model.ID),
		Emoji:  model.Emoji,
		UserID: model.Edges.Account.ID.String(),
		PostID: model.Edges.Post.ID.String(),
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
