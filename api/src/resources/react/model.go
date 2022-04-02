package react

import (
	"github.com/forPelevin/gomoji"
	"github.com/Southclaws/storyden/api/src/infra/db"
)

type React struct {
	ID    string `json:"id"`
	Emoji string `json:"emoji"`
	User  string `json:"user"`
	Post  string `json:"post"`
}

func FromModel(model *db.ReactModel, postID string) *React {
	return &React{
		ID:    model.ID,
		Emoji: model.Emoji,
		User:  model.RelationsReact.User.ID,
		Post:  postID,
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
