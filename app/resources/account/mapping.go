package account

import (
	"net/mail"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

func MapRef(a *ent.Account) (*Account, error) {
	bio, err := datagraph.NewRichText(a.Bio)
	if err != nil {
		return nil, err
	}

	kind, err := NewAccountKind(string(a.Kind))
	if err != nil {
		return nil, err
	}

	var roles held.Roles
	if rolesEdge := a.Edges.AccountRoles; rolesEdge != nil {
		roles, err = held.MapList(rolesEdge)
		if err != nil {
			return nil, err
		}
	}

	return &Account{
		ID:        AccountID(a.ID),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,

		Handle:   a.Handle,
		Name:     a.Name,
		Bio:      bio,
		Kind:     kind,
		Admin:    a.Admin, // TODO: should this be derived from roles?
		Roles:    roles,
		Metadata: a.Metadata,

		DeletedAt: opt.NewPtr(a.DeletedAt),
		IndexedAt: opt.NewPtr(a.IndexedAt),
	}, nil
}

func MapAccount(a *ent.Account) (*AccountWithEdges, error) {
	ref, err := MapRef(a)
	if err != nil {
		return nil, err
	}

	authsEdge, err := a.Edges.AuthenticationOrErr()
	if err != nil {
		return nil, err
	}

	emailsEdge, err := a.Edges.EmailsOrErr()
	if err != nil {
		return nil, err
	}

	auths := dt.Map(authsEdge, func(a *ent.Authentication) string {
		return a.Service
	})

	verifiedStatus := VerifiedStatusNone
	if len(dt.Filter(emailsEdge, func(e *ent.Email) bool { return e.Verified })) > 0 {
		verifiedStatus = VerifiedStatusVerifiedEmail
	}

	emails := dt.Map(emailsEdge, MapEmail)

	invitedByEdge := opt.NewPtr(a.Edges.InvitedBy)

	invitedBy, err := opt.MapErr(invitedByEdge, func(i ent.Invitation) (Account, error) {
		c, err := i.Edges.CreatorOrErr()
		if err != nil {
			return Account{}, err
		}

		ib, err := MapRef(c)
		if err != nil {
			return Account{}, err
		}

		return *ib, nil
	})
	if err != nil {
		return nil, err
	}

	links, err := dt.MapErr(a.Links, MapExternalLink)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &AccountWithEdges{
		Account:        *ref,
		Auths:          auths,
		EmailAddresses: emails,
		VerifiedStatus: verifiedStatus,
		InvitedBy:      invitedBy,
		ExternalLinks:  links,
	}, nil
}

func MapEmail(in *ent.Email) *EmailAddress {
	addr, _ := mail.ParseAddress(in.EmailAddress)
	// NOTE: Ent already validates this
	// TODO: use mail.Address instead of string in ent schema

	return &EmailAddress{
		ID:       in.ID,
		Email:    *addr,
		Verified: in.Verified,
	}
}
