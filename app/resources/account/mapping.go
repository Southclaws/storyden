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

	return &Account{
		ID:        AccountID(a.ID),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,

		Handle:   a.Handle,
		Name:     a.Name,
		Bio:      bio,
		Kind:     kind,
		Admin:    a.Admin, // TODO: should this be derived from roles?
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

	rolesEdge := a.Edges.AccountRoles

	auths := dt.Map(a.Edges.Authentication, func(a *ent.Authentication) string {
		return a.Service
	})

	roles, err := held.MapList(rolesEdge, a.Admin)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	verifiedStatus := VerifiedStatusNone
	if len(dt.Filter(a.Edges.Emails, func(e *ent.Email) bool { return e.Verified })) > 0 {
		verifiedStatus = VerifiedStatusVerifiedEmail
	}

	emails := dt.Map(a.Edges.Emails, MapEmail)

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

	return &AccountWithEdges{
		Account:        *ref,
		Roles:          roles,
		Auths:          auths,
		EmailAddresses: emails,
		VerifiedStatus: verifiedStatus,
		InvitedBy:      invitedBy,
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
