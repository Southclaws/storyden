package account

import (
	"net/mail"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

func MapAccount(a *ent.Account) (*Account, error) {
	auths := dt.Map(a.Edges.Authentication, func(a *ent.Authentication) string {
		return a.Service
	})

	bio, err := content.NewRichText(a.Bio)
	if err != nil {
		return nil, err
	}

	links, err := dt.MapErr(a.Links, MapExternalLink)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	verifiedStatus := VerifiedStatusNone
	if len(dt.Filter(a.Edges.Emails, func(e *ent.Email) bool { return e.Verified })) > 0 {
		verifiedStatus = VerifiedStatusVerifiedEmail
	}

	emails := dt.Map(a.Edges.Emails, MapEmail)

	return &Account{
		ID:             AccountID(a.ID),
		Handle:         a.Handle,
		Name:           a.Name,
		Bio:            bio,
		Admin:          a.Admin,
		Auths:          auths,
		EmailAddresses: emails,
		VerifiedStatus: verifiedStatus,
		ExternalLinks:  links,
		Metadata:       a.Metadata,

		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: opt.NewPtr(a.DeletedAt),
	}, nil
}

func MapExternalLink(e schema.ExternalLink) (ExternalLink, error) {
	u, err := url.Parse(e.URL)
	if err != nil {
		return ExternalLink{}, err
	}

	return ExternalLink{
		Text: e.Text,
		URL:  *u,
	}, nil
}

func MapEmail(in *ent.Email) *EmailAddress {
	addr, _ := mail.ParseAddress(in.EmailAddress)
	// NOTE: Ent already validates this
	// TODO: use mail.Address instead of string in ent schema

	return &EmailAddress{
		Email:    *addr,
		Verified: in.Verified,
		IsAuth:   in.AuthenticationRecordID != nil,
	}
}
