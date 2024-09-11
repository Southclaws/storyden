package profile

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/internal/ent"
)

type Public struct {
	ID      account.AccountID
	Created time.Time
	Deleted opt.Optional[time.Time]

	Handle        string
	Name          string
	Bio           content.Rich
	Admin         bool
	Followers     int
	Following     int
	LikeScore     int
	Interests     []*tag.Tag
	ExternalLinks []account.ExternalLink
	Metadata      map[string]any
}

func (p *Public) GetID() xid.ID             { return xid.ID(p.ID) }
func (p *Public) GetKind() datagraph.Kind   { return datagraph.KindProfile }
func (p *Public) GetName() string           { return p.Name }
func (p *Public) GetSlug() string           { return p.Handle }
func (p *Public) GetDesc() string           { return p.Bio.Short() }
func (p *Public) GetContent() content.Rich  { return p.Bio }
func (p *Public) GetProps() map[string]any  { return p.Metadata }
func (p *Public) GetAssets() []*asset.Asset { return []*asset.Asset{} }

func ProfileFromModel(a *ent.Account) (*Public, error) {
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

	return &Public{
		ID:        account.AccountID(a.ID),
		Created:   a.CreatedAt,
		Deleted:   opt.NewPtr(a.DeletedAt),
		Handle:    a.Handle,
		Name:      a.Name,
		Bio:       bio,
		Admin:     a.Admin,
		Interests: interests,
		Metadata:  a.Metadata,
	}, nil
}

func ProfileFromAccount(a *account.Account) *Public {
	return &Public{
		ID:            a.ID,
		Created:       a.CreatedAt,
		Deleted:       a.DeletedAt,
		Handle:        a.Handle,
		Name:          a.Name,
		Bio:           a.Bio,
		Admin:         a.Admin,
		Followers:     a.Followers,
		Following:     a.Following,
		LikeScore:     a.LikeScore,
		Interests:     nil,
		ExternalLinks: a.ExternalLinks,
		Metadata:      a.Metadata,
	}
}
