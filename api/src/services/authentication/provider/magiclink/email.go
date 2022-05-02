package magiclink

import (
	"github.com/Southclaws/storyden/api/src/infra/mailer"
	"github.com/Southclaws/storyden/api/src/resources/authentication"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

type Email struct {
	repo authentication.Repository
	m    mailer.Mailer
}

func NewEmail(repo authentication.Repository, m mailer.Mailer) *Email {
	return &Email{
		repo,
		m,
	}
}

// Send sends a magic link
func (a *Email) Send(email string) (*user.User, error) {
	return nil, nil
}

func (a *Email) Callback(token []byte) (*user.User, error) {
	return nil, nil
}
