package magiclink

import (
	"github.com/Southclaws/storyden/backend/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
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
func (a *Email) Send(email string) (*account.Account, error) {
	return nil, nil
}

func (a *Email) Callback(token []byte) (*account.Account, error) {
	return nil, nil
}
