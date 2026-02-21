package plugin_runner

import (
	"context"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/duplex"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

// Host manages all plugins, both supervised and external, and their sessions.
// Supervised plugins are plugins which are installed via .zip files and run as
// child processes of the Storyden instance. External plugins are plugins which
// connect to the Storyden instance externally via WebSocket and authenticate
// in a slightly different way. Their lifecycle is not controlled by the Host,
// hosever the Host has authority on disconnecting them if they are unloaded.
type Host interface {
	// Connects a supervised or external plugin to the host.
	Connect(ctx context.Context, id plugin.InstallationID, duplex duplex.Duplex) error

	// Load loads a plugin into the host.
	Load(ctx context.Context, pr plugin.Record) error

	// Unloads a plugin, if supervised, stops the process.
	Unload(ctx context.Context, id plugin.InstallationID) error

	// Returns the session by ID.
	GetSession(ctx context.Context, id plugin.InstallationID) (Session, error)

	// Lists all connected plugins, both supervised and external.
	GetSessions(ctx context.Context) ([]Session, error)
}

// Supervised represents an in process plugin which is run as a child process of
// the Storyden instance. The Host manages the lifecycle of supervised plugins,
// starting and stopping them as needed.
type Supervised interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Session represents any connected plugin, both supervised and external. It
// provides methods for getting its status, sending RPCs and the initial Connect
// handshake which attaches the Duplex connection to the Session. The Host does
// not manage the lifecycle of sessions, they are responsible for managing their
// own lifecycle. If a Session represents a supervised plugin, it must handle
// its own crash handling, restarting and logging.
type Session interface {
	// ID returns the plugin installation ID of the session's plugin.
	ID() plugin.InstallationID

	// Supervised returns nil if the session represents an external plugin.
	Supervised() Supervised

	// Not necessarily needed?
	GetStartedAt() opt.Optional[time.Time]
	GetReportedState() plugin.ReportedState
	GetErrorMessage() string

	// SetActiveState changes the desired active state of the plugin.
	// For supervised plugins, this starts/stops the underlying process.
	// For external plugins, Active->Inactive disconnects the websocket session.
	SetActiveState(ctx context.Context, state plugin.ActiveState) error

	// Connect is called by the Host when a plugin connects. This applies to
	// both supervised and external plugins. The Host determines whether the
	// authentication information is valid for either type and passes the
	// connection to the Session.
	Connect(ctx context.Context, duplex duplex.Duplex) error

	// Send sends an RPC command to a connected plugin and returns the response.
	Send(ctx context.Context, id xid.ID, payload rpc.HostToPluginRequestUnion) (rpc.HostToPluginResponseUnion, error)
}
