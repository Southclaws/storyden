package supervised_runtime

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Southclaws/storyden/app/resources/plugin"
)

type Event struct {
	State   plugin.ReportedState
	Message string
}

type Runtime interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Events() <-chan Event
}

type Provider interface {
	New(
		id plugin.InstallationID,
		bin []byte,
		manifest *plugin.Validated,
		parentLogger *slog.Logger,
		serverURL url.URL,
		parentCtx context.Context,
	) Runtime
}
