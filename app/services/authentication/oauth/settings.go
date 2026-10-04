package oauth

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func (s *Service) registrationSettings(ctx context.Context) (settings.OAuthServiceSettings, error) {
	configuration, err := s.settings.Get(ctx)
	if err != nil {
		return settings.OAuthServiceSettings{}, err
	}
	return configuration.Services.OrZero().OAuth.OrZero(), nil
}
