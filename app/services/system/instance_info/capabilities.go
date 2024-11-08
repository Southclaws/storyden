package instance_info

//go:generate go run github.com/Southclaws/enumerator

type capabilityEnum string

const (
	capabilityNone   capabilityEnum = `none`
	capabilitySemdex capabilityEnum = `semdex`
)

type Capabilities []Capability
