package plugin

//go:generate go run github.com/Southclaws/enumerator

type capabilityEnum string

const (
	capabilityNetwork  capabilityEnum = "network"
	capabilityDatabase capabilityEnum = "database"
)
