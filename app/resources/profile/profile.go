package profile

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/asset"

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
	Bio           datagraph.Content
	Admin         bool
	Followers     int
	Following     int
	LikeScore     int
	Roles         held.Roles
	Interests     []*tag.Tag
	ExternalLinks []account.ExternalLink
	InvitedBy     opt.Optional[Public]
	Metadata      map[string]any
}

func (p *Public) GetID() xid.ID                 { return xid.ID(p.ID) }
func (p *Public) GetKind() datagraph.Kind       { return datagraph.KindProfile }
func (p *Public) GetName() string               { return p.Name }
func (p *Public) GetSlug() string               { return p.Handle }
func (p *Public) GetDesc() string               { return p.Bio.Short() }
func (p *Public) GetContent() datagraph.Content { return p.Bio }
func (p *Public) GetProps() map[string]any      { return p.Metadata }
func (p *Public) GetAssets() []*asset.Asset     { return []*asset.Asset{} }
func (p *Public) GetCreated() time.Time         { return p.Created }
func (p *Public) GetUpdated() time.Time         { return p.Created }

func ProfileFromModel(a *ent.Account) (*Public, error) {
	rolesEdge, err := a.Edges.AccountRolesOrErr()
	if err != nil {
		return nil, err
	}

	roles, err := held.MapList(rolesEdge, a.Admin)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	interests := dt.Map(a.Edges.Tags, func(t *ent.Tag) *tag.Tag {
		return &tag.Tag{
			ID:   t.ID.String(),
			Name: t.Name,
		}
	})

	bio, err := datagraph.NewRichText(a.Bio)
	if err != nil {
		return nil, err
	}

	invitedByEdge := opt.NewPtr(a.Edges.InvitedBy)

	invitedBy, err := opt.MapErr(invitedByEdge, func(i ent.Invitation) (Public, error) {
		c, err := i.Edges.CreatorOrErr()
		if err != nil {
			return Public{}, err
		}

		ib, err := account.MapAccount(c)
		if err != nil {
			return Public{}, err
		}

		return *ProfileFromAccount(ib), nil
	})
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
		Roles:     roles,
		Interests: interests,
		InvitedBy: invitedBy,
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
		Roles:         a.Roles,
		Interests:     nil,
		ExternalLinks: a.ExternalLinks,
		InvitedBy: opt.Map(a.InvitedBy, func(a account.Account) Public {
			return *ProfileFromAccount(&a)
		}),
		Metadata: a.Metadata,
	}
}
