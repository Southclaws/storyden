package robot

import (
	"encoding/json"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
)

type TurnID xid.ID

func (id TurnID) String() string { return xid.ID(id).String() }

type TurnStatus string

const (
	SessionEventPageSize = 256

	TurnStatusQueued    TurnStatus = "queued"
	TurnStatusRunning   TurnStatus = "running"
	TurnStatusCompleted TurnStatus = "completed"
	TurnStatusBlocked   TurnStatus = "blocked"
	TurnStatusFailed    TurnStatus = "failed"
	TurnStatusCancelled TurnStatus = "cancelled"
)

type Turn struct {
	ID                TurnID
	SessionID         SessionID
	InitiatorID       opt.Optional[account.AccountID]
	SourceKind        string
	RobotRef          string
	InputData         json.RawMessage
	Status            TurnStatus
	CreatedAt         time.Time
	StartedAt         opt.Optional[time.Time]
	FinishedAt        opt.Optional[time.Time]
	CancelRequestedAt opt.Optional[time.Time]
	ErrorText         opt.Optional[string]
}

type SessionEventKind string

const (
	SessionEventMessage       SessionEventKind = "message"
	SessionEventInputQueued   SessionEventKind = "input_queued"
	SessionEventTurnQueued    SessionEventKind = "turn_queued"
	SessionEventTurnCompleted SessionEventKind = "turn_completed"
	SessionEventTurnBlocked   SessionEventKind = "turn_blocked"
	SessionEventTurnFailed    SessionEventKind = "turn_failed"
	SessionEventTurnCancelled SessionEventKind = "turn_cancelled"
)

type SessionEvent struct {
	Sequence  uint64
	Kind      SessionEventKind
	TurnID    opt.Optional[TurnID]
	InputIDs  []InputID
	Message   opt.Optional[*Message]
	ErrorText opt.Optional[string]
}
