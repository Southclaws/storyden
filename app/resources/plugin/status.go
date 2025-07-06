package plugin

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/internal/ent"
)

//go:generate go run github.com/Southclaws/enumerator

type activeStateEnum string

const (
	activeStateActive   activeStateEnum = "active"
	activeStateInactive activeStateEnum = "inactive"
	activeStateError    activeStateEnum = "error"
)

type Status struct {
	ActiveState   ActiveState
	ChangedAt     time.Time
	StatusMessage string
	Details       map[string]any
}

func MapStatus(in *ent.Plugin) (*Status, error) {
	activeState, err := NewActiveState(in.ActiveState)
	if err != nil {
		return nil, err
	}

	return &Status{
		ActiveState:   activeState,
		ChangedAt:     in.ActiveStateChangedAt,
		StatusMessage: opt.NewPtr(in.StatusMessage).OrZero(),
		Details:       in.StatusDetails,
	}, nil
}
