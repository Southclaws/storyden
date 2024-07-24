package account

//go:generate go run github.com/Southclaws/enumerator

type verifiedStatusEnum string

const (
	verifiedStatusNone          verifiedStatusEnum = "none"
	verifiedStatusVerifiedEmail verifiedStatusEnum = "verified_email"
)
