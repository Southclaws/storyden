package account

import (
	"net/mail"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

func mapRefWithoutRoles(accID xid.ID) (held.Roles, error) {
	return held.Roles{}, nil
}

func MapRef(a *ent.Account) (*Account, error) {
	return RefMapper(mapRefWithoutRoles)(a)
}

func RefMapper(roleHydratorFn func(accID xid.ID) (held.Roles, error)) func(a *ent.Account) (*Account, error) {
	return func(a *ent.Account) (*Account, error) {
		bio, err := datagraph.NewRichText(a.Bio)
		if err != nil {
			return nil, err
		}

		kind, err := NewAccountKind(string(a.Kind))
		if err != nil {
			return nil, err
		}

		roles, err := roleHydratorFn(a.ID)
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
			Roles:    roles,
			Admin:    a.Admin, // TODO: should this be derived from roles?
			Metadata: a.Metadata,

			DeletedAt: opt.NewPtr(a.DeletedAt),
			IndexedAt: opt.NewPtr(a.IndexedAt),
		}, nil
	}
}

func MapAccount(roleHydratorFn func(accID xid.ID) (held.Roles, error)) func(a *ent.Account) (*AccountWithEdges, error) {
	refMapper := RefMapper(roleHydratorFn)

	return func(a *ent.Account) (*AccountWithEdges, error) {
		ref, err := refMapper(a)
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

			ib, err := refMapper(c)
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
			Roles:          ref.Roles,
			Auths:          auths,
			EmailAddresses: emails,
			VerifiedStatus: verifiedStatus,
			InvitedBy:      invitedBy,
			ExternalLinks:  links,
		}, nil
	}
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
