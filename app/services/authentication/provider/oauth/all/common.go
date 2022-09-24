// Package all contains things that are common amongst all providers.
package all

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
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
func LoadProvider(name string) (Configuration, error) {
	enabled := struct{ Enabled bool }{}
	if envconfig.Process(strings.ToUpper(name), &enabled); !enabled.Enabled {
		return Configuration{}, nil
	}

	pc := Configuration{}
	if err := envconfig.Process(strings.ToUpper(name), &pc); err != nil {
		return Configuration{}, errors.Wrapf(err, "oauth provider '%s' is enabled but configuration failed to load", name)
	}

	return pc, nil
}
