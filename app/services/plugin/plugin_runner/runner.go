package plugin_runner

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/lib/plugin"
)

type Runner interface {
	Load(ctx context.Context, bin []byte) (*PluginSession, error)
	Unload(ctx context.Context, id plugin.ID) error
	Validate(ctx context.Context, bin []byte) (*plugin.Manifest, error)
	GetSession(ctx context.Context, id plugin.ID) (*PluginSession, error)
	GetSessions(ctx context.Context) ([]*PluginSession, error)

	StartPlugin(ctx context.Context, id plugin.ID) error
	StopPlugin(ctx context.Context, id plugin.ID) error
}

func Build() fx.Option {
	return fx.Provide(
	// TODO: Provide a runner implementation to the system.
	)
}
