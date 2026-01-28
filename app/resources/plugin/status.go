package plugin

//go:generate go run github.com/Southclaws/enumerator

type activeStateEnum string

const (
	activeStateActive   activeStateEnum = "active"
	activeStateInactive activeStateEnum = "inactive"
)

type reportedStateEnum string

const (
	reportedStateActive   reportedStateEnum = "active"
	reportedStateInactive reportedStateEnum = "inactive"
	reportedStateError    reportedStateEnum = "errored"
)
