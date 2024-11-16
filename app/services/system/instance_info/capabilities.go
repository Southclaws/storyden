package instance_info

//go:generate go run github.com/Southclaws/enumerator

type capabilityEnum string

const (
	capabilitySemdex      capabilityEnum = `semdex`
	capabilityEmailClient capabilityEnum = `email_client`
	capabilitySMSClient   capabilityEnum = `sms_client`
)

type Capabilities []Capability
