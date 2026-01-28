package plugin_runner

import (
	"context"
	"time"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/plugin"
)

type Runner interface {
	Load(ctx context.Context, id plugin.InstallationID, bin []byte) (Session, error)
	Unload(ctx context.Context, id plugin.InstallationID) error

	Validate(ctx context.Context, bin []byte) (*plugin.Validated, error)

	GetSession(ctx context.Context, id plugin.InstallationID) (Session, error)
	GetSessions(ctx context.Context) ([]Session, error)
}

type Session interface {
	ID() plugin.InstallationID

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	GetStartedAt() opt.Optional[time.Time]
	GetReportedState() plugin.ReportedState
	GetErrorMessage() string

	Send(ctx context.Context, method string, params any) (any, error)
	Connect(ctx context.Context, duplex Duplex) error
}

type Duplex interface {
	Send(ctx context.Context, b []byte) error
	Recv(ctx context.Context) ([]byte, error)
	Close() error
}
