// Package all contains things that are common amongst all providers.
package all

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Configuration is the standard config required by all OAuth2 providers.
type Configuration struct {
	Enabled      bool   `default:"false"`
	ClientID     string `required:"true"`
	ClientSecret string `required:"true"`
}

// LoadProvider attempts to load a configuration for an OAuth2 provider from
// environment variables. The way this works is, if the parse fails then the
// provider is considered disabled and an empty configuration is returned.
func LoadProvider(name string) Configuration {
	pc := Configuration{}
	if err := envconfig.Process(strings.ToUpper(name), &pc); err != nil {
		return Configuration{}
	}

	return pc
}
