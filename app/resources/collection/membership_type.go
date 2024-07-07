package collection

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type membershipTypeEnum string

const (
	membershipTypeNormal     membershipTypeEnum = "normal"
	membershipTypeSubmission membershipTypeEnum = "submission"
)
