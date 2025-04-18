package session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/token"
)

type Issuer struct {
	tokenRepo token.Repository
}

func NewIssuer(tokenRepo token.Repository) *Issuer {
	return &Issuer{
		tokenRepo: tokenRepo,
	}
}

func (s *Issuer) Issue(ctx context.Context, accountID account.AccountID) (*token.Token, error) {
	t, err := s.tokenRepo.Issue(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &t.Token, nil
}
