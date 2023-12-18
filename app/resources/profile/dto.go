package profile

import (
	"time"

	"github.com/Southclaws/dt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/ent"
)

type Profile struct {
	ID      account.AccountID
	Created time.Time

	Handle    string
	Name      string
	Bio       string
	Admin     bool
	Interests []*tag.Tag
}

func FromModel(a *ent.Account) (*Profile, error) {
	interests := dt.Map(a.Edges.Tags, func(t *ent.Tag) *tag.Tag {
		return &tag.Tag{
			ID:   t.ID.String(),
			Name: t.Name,
		}
	})

	return &Profile{
		ID:        account.AccountID(a.ID),
		Created:   a.CreatedAt,
		Handle:    a.Handle,
		Name:      a.Name,
		Bio:       a.Bio,
		Admin:     a.Admin,
		Interests: interests,
	}, nil
}
