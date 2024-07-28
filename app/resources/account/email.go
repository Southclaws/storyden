package account

import (
	"net/mail"
)

type EmailAddress struct {
	Email    mail.Address
	Verified bool
	IsAuth   bool
}
