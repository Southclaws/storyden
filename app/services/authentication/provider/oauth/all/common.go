// Package all contains things that are common amongst all providers.
package all

import (
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/internal/config"
)

// Configuration is the standard config required by all OAuth2 providers.
type Configuration struct {
	Enabled      bool   `default:"false"`
	ClientID     string `required:"true" split_words:"true"`
	ClientSecret string `required:"true" split_words:"true"`
}

// LoadProvider attempts to load a configuration for an OAuth2 provider from
// environment variables. The way this works is, if the parse fails then the
// provider is considered disabled and an empty configuration is returned.
func LoadProvider(svc authentication.Service) (*Configuration, error) {
	name := svc.String()

	enabled := struct{ Enabled bool }{}
	if envconfig.Process(strings.ToUpper(name), &enabled); !enabled.Enabled {
		return nil, nil
	}

	pc := Configuration{}
	if err := envconfig.Process(strings.ToUpper(name), &pc); err != nil {
		return nil, fault.Wrap(err, fmsg.With(fmt.Sprintf("oauth provider '%s' is enabled but configuration failed to load", name)))
	}

	return &pc, nil
}

// TODO: Let the frontend construct this URL via some kind of template.
func Redirect(cfg config.Config, svc authentication.Service) string {
	name := svc.String()

	return fmt.Sprintf("%s/auth/%s/callback", cfg.PublicWebAddress.String(), name)
}
