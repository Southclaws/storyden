package plugin

//go:generate go run github.com/Southclaws/enumerator

type activeStateEnum string

const (
	activeStateActive   activeStateEnum = "active"
	activeStateInactive activeStateEnum = "inactive"
	activeStateError    activeStateEnum = "error"
)
