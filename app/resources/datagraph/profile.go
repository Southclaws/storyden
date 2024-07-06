package datagraph

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/ent"
)

type Profile struct {
	ID      account.AccountID
	Created time.Time
	Deleted opt.Optional[time.Time]

	Handle        string
	Name          string
	Bio           content.Rich
	Admin         bool
	Interests     []*tag.Tag
	ExternalLinks []account.ExternalLink
}

func (p *Profile) GetID() xid.ID   { return xid.ID(p.ID) }
func (p *Profile) GetKind() Kind   { return KindProfile }
func (p *Profile) GetName() string { return p.Name }
func (p *Profile) GetSlug() string { return p.Handle }
func (p *Profile) GetDesc() string { return p.Bio.Short() }
func (p *Profile) GetText() string { return p.Bio.HTML() }
func (p *Profile) GetProps() any   { return nil }

func ProfileFromModel(a *ent.Account) (*Profile, error) {
	interests := dt.Map(a.Edges.Tags, func(t *ent.Tag) *tag.Tag {
		return &tag.Tag{
			ID:   t.ID.String(),
			Name: t.Name,
		}
	})

	bio, err := content.NewRichText(a.Bio)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:        account.AccountID(a.ID),
		Created:   a.CreatedAt,
		Deleted:   opt.NewPtr(a.DeletedAt),
		Handle:    a.Handle,
		Name:      a.Name,
		Bio:       bio,
		Admin:     a.Admin,
		Interests: interests,
	}, nil
}

func ProfileFromAccount(a *account.Account) *Profile {
	return &Profile{
		ID:            a.ID,
		Created:       a.CreatedAt,
		Deleted:       a.DeletedAt,
		Handle:        a.Handle,
		Name:          a.Name,
		Bio:           a.Bio,
		Admin:         a.Admin,
		Interests:     nil,
		ExternalLinks: a.ExternalLinks,
	}
}
