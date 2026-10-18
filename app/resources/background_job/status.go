package background_job

//go:generate go run github.com/Southclaws/enumerator

type statusEnum string

const (
	statusPending    statusEnum = "pending"
	statusProcessing statusEnum = "processing"
	statusSkipped    statusEnum = "skipped"
	statusFailed     statusEnum = "failed"
	statusComplete   statusEnum = "complete"
)
