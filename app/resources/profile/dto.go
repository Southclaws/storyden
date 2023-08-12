package profile

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/ent"
)

type Profile struct {
	AccountID account.AccountID

	Interests []*tag.Tag
}

func FromModel(a *ent.Account) (*Profile, error) {
	tags, err := a.Edges.TagsOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	interests := dt.Map(tags, func(t *ent.Tag) *tag.Tag {
		return &tag.Tag{
			ID:   t.ID.String(),
			Name: t.Name,
		}
	})

	return &Profile{
		Interests: interests,
	}, nil
}
