package email_queue

import "github.com/Southclaws/storyden/internal/ent/emailqueue"

//go:generate go run github.com/Southclaws/enumerator

type statusEnum string

const (
	statusPending    statusEnum = "pending"
	statusProcessing statusEnum = "processing"
	statusSent       statusEnum = "sent"
	statusFailed     statusEnum = "failed"
)

func (s Status) Ent() emailqueue.Status {
	return emailqueue.Status(s.String())
}
