package account

import (
	"net/mail"

	"github.com/rs/xid"
)

type EmailAddress struct {
	ID       xid.ID
	Email    mail.Address
	Verified bool
}
