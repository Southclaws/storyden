package background_job

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type State struct {
	ID            xid.ID
	Kind          string
	Key           string
	Revision      int64
	Status        Status
	Attempts      int
	Reason        string
	LeaseToken    string
	LeaseUntil    time.Time
	DispatchUntil time.Time
	AvailableAt   time.Time
}

type Job[T any] struct {
	State
	Payload T
}

func Map[T any](row *ent.BackgroundJob) (*Job[T], error) {
	status, err := NewStatus(string(row.Status))
	if err != nil {
		return nil, err
	}
	payload, err := Decode[T](row.Payload)
	if err != nil {
		return nil, err
	}

	return &Job[T]{
		State: State{
			ID:            row.ID,
			Kind:          row.Kind,
			Key:           row.Key,
			Revision:      row.Revision,
			Status:        status,
			Attempts:      row.Attempts,
			Reason:        row.Reason,
			LeaseToken:    row.LeaseToken,
			LeaseUntil:    row.LeaseUntil,
			DispatchUntil: row.DispatchUntil,
			AvailableAt:   row.AvailableAt,
		},
		Payload: payload,
	}, nil
}

func Decode[T any](raw json.RawMessage) (T, error) {
	var payload T
	err := json.Unmarshal(raw, &payload)
	return payload, err
}
