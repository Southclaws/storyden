package oauth

import (
	"context"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
)

type UserInfo struct {
	Subject           string
	Name              opt.Optional[string]
	Email             opt.Optional[string]
	EmailVerified     opt.Optional[bool]
	PreferredUsername opt.Optional[string]
}

func (s *Service) UserInfo(ctx context.Context, accountID account.AccountID, scopes []string) (*UserInfo, error) {
	out := &UserInfo{
		Subject: accountID.String(),
	}

	if !contains(scopes, "profile") && !contains(scopes, "email") {
		return out, nil
	}

	acc, err := s.account.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if contains(scopes, "profile") {
		out.Name = opt.New(acc.Name)
		out.PreferredUsername = opt.New(acc.Handle)
	}

	if contains(scopes, "email") {
		email := ""
		emailVerified := false
		for _, address := range acc.EmailAddresses {
			email = address.Email.Address
			emailVerified = address.Verified
			if address.Verified {
				break
			}
		}
		out.Email = opt.New(email)
		out.EmailVerified = opt.New(emailVerified)
	}

	return out, nil
}
