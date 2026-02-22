package plugin_runner

import (
	"strings"

	"github.com/Southclaws/fault"
)

type RuntimeProvider string

const (
	RuntimeProviderLocal   RuntimeProvider = "local"
	RuntimeProviderSprites RuntimeProvider = "sprites"
)

func (p RuntimeProvider) String() string {
	return string(p)
}

func ParseRuntimeProvider(v string) (RuntimeProvider, error) {
	value := strings.TrimSpace(strings.ToLower(v))
	if value == "" {
		return RuntimeProviderLocal, nil
	}

	provider := RuntimeProvider(value)
	switch provider {
	case RuntimeProviderLocal, RuntimeProviderSprites:
		return provider, nil
	default:
		return "", fault.Newf("unknown plugin runtime provider: %q", v)
	}
}
