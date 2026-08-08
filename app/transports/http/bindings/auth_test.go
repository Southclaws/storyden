package bindings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

type fakeProvider struct {
	service authentication.Service
}

func (f fakeProvider) Enabled(ctx context.Context) (bool, error) { return true, nil }
func (f fakeProvider) Service() authentication.Service           { return f.service }
func (f fakeProvider) Token() authentication.TokenType           { return authentication.TokenTypeOAuth }

type fakeDisplayNameProvider struct {
	fakeProvider
	name string
}

func (f fakeDisplayNameProvider) DisplayName() string { return f.name }

func TestProviderNameFallsBackToServiceFormat(t *testing.T) {
	t.Parallel()

	p := fakeProvider{service: authentication.ServiceOAuthGitHub}

	assert.Equal(t, "GitHub", providerName(p))
}

func TestProviderNameUsesDisplayNamerOverride(t *testing.T) {
	t.Parallel()

	p := fakeDisplayNameProvider{
		fakeProvider: fakeProvider{service: authentication.ServiceOAuthKeycloak},
		name:         "Corporate SSO",
	}

	assert.Equal(t, "Corporate SSO", providerName(p))
}

func TestProviderNameFallsBackWhenDisplayNameEmpty(t *testing.T) {
	t.Parallel()

	p := fakeDisplayNameProvider{
		fakeProvider: fakeProvider{service: authentication.ServiceOAuthKeycloak},
		name:         "",
	}

	assert.Equal(t, "Keycloak", providerName(p))
}
