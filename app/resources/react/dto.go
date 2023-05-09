package react

import (
	"github.com/forPelevin/gomoji"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type ReactID xid.ID

type React struct {
	ID     ReactID
	Emoji  string
	UserID string
	PostID string
}

func FromModel(in *ent.React) *React {
	return &React{
		ID:     ReactID(in.ID),
		Emoji:  in.Emoji,
		UserID: in.AccountID.String(),
		PostID: in.PostID.String(),
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
