package account

import (
	"net/mail"

	"github.com/Southclaws/storyden/internal/ent"
)

type EmailAddress struct {
	Email    mail.Address
	Verified bool
	IsAuth   bool
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
