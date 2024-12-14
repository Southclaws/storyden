package instance_info

//go:generate go run github.com/Southclaws/enumerator

type capabilityEnum string

const (
	capabilityGenAI       capabilityEnum = `gen_ai`
	capabilitySemdex      capabilityEnum = `semdex`
	capabilityEmailClient capabilityEnum = `email_client`
	capabilitySMSClient   capabilityEnum = `sms_client`
)

type Capabilities []Capability

func (c Capabilities) Has(capability Capability) bool {
	for _, cap := range c {
		if cap == capability {
			return true
		}
	}
	return false
}
