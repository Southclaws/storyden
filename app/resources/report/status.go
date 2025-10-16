package report

//go:generate go run github.com/Southclaws/enumerator

type statusEnum string

const (
	statusSubmitted    statusEnum = "submitted"
	statusAcknowledged statusEnum = "acknowledged"
	statusResolved     statusEnum = "resolved"
)
