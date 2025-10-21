package plugin

//go:generate go run github.com/Southclaws/enumerator

type activeStateEnum string

const (
	activeStateActive   activeStateEnum = "active"
	activeStateInactive activeStateEnum = "inactive"
)

type reportedStateEnum string

const (
	reportedStateInactive   reportedStateEnum = "inactive"
	reportedStateStarting   reportedStateEnum = "starting"
	reportedStateConnecting reportedStateEnum = "connecting"
	reportedStateActive     reportedStateEnum = "active"
	reportedStateStopping   reportedStateEnum = "stopping"
	reportedStateError      reportedStateEnum = "errored"
	reportedStateRestarting reportedStateEnum = "restarting"
)

type pluginModeEnum string

const (
	pluginModeSupervised pluginModeEnum = "supervised"
	pluginModeExternal   pluginModeEnum = "external"
)
