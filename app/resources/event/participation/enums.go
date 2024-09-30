package participation

//go:generate go run github.com/Southclaws/enumerator

type policyEnum string

const (
	policyOpen       policyEnum = "open"
	policyClosed     policyEnum = "closed"
	policyInviteOnly policyEnum = "invite_only"
)

type roleEnum string

const (
	roleHost     roleEnum = "host"
	roleAttendee roleEnum = "attendee"
)

type statusEnum string

const (
	statusRequested statusEnum = "requested"
	statusInvited   statusEnum = "invited"
	statusAttending statusEnum = "attending"
	statusDeclined  statusEnum = "declined"
)
