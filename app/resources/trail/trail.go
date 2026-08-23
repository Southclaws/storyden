package trail

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
)

//go:generate go run github.com/Southclaws/enumerator

type (
	ID          xid.ID
	ActionID    xid.ID
	RunID       xid.ID
	ActionRunID xid.ID
)

func (id ID) String() string          { return xid.ID(id).String() }
func (id ActionID) String() string    { return xid.ID(id).String() }
func (id RunID) String() string       { return xid.ID(id).String() }
func (id ActionRunID) String() string { return xid.ID(id).String() }

type statusEnum string

const (
	statusActive   statusEnum = "active"
	statusPaused   statusEnum = "paused"
	statusFinished statusEnum = "finished"
	statusArchived statusEnum = "archived"
)

type triggerTypeEnum string

type runKindEnum string

const (
	runKindScheduled runKindEnum = "scheduled"
	runKindEvent     runKindEnum = "event"
	runKindManual    runKindEnum = "manual"
)

type runStatusEnum string

const (
	runStatusQueued            runStatusEnum = "queued"
	runStatusRunning           runStatusEnum = "running"
	runStatusCompleted         runStatusEnum = "completed"
	runStatusAttentionRequired runStatusEnum = "attention_required"
	runStatusCancelled         runStatusEnum = "cancelled"
	runStatusSkipped           runStatusEnum = "skipped"
)

type actionRunStatusEnum string

const (
	actionRunStatusQueued    actionRunStatusEnum = "queued"
	actionRunStatusRunning   actionRunStatusEnum = "running"
	actionRunStatusCompleted actionRunStatusEnum = "completed"
	actionRunStatusBlocked   actionRunStatusEnum = "blocked"
	actionRunStatusFailed    actionRunStatusEnum = "failed"
	actionRunStatusCancelled actionRunStatusEnum = "cancelled"
)

type actionKindEnum string

const (
	triggerTypeSchedule triggerTypeEnum = "schedule"
	triggerTypeEvent    triggerTypeEnum = "event"
	actionKindRobotRun  actionKindEnum  = "robot_run"
)

type robotInvocationOutputStatusEnum string

const (
	robotInvocationOutputStatusCompleted robotInvocationOutputStatusEnum = "completed"
	robotInvocationOutputStatusBlocked   robotInvocationOutputStatusEnum = "blocked"
	robotInvocationOutputStatusFailed    robotInvocationOutputStatusEnum = "failed"
)

type Trail struct {
	ID               ID
	Creator          account.Account
	Name             string
	Description      string
	Status           Status
	Trigger          Trigger
	NextOccurrenceAt *time.Time
	LastOccurrenceAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Actions          []*Action
}

type Action struct {
	ID         ActionID
	TrailID    ID
	Kind       ActionKind
	Position   int
	Config     json.RawMessage
	ArchivedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TriggerEvent struct {
	TrailID      string
	TrailRunID   string
	Kind         RunKind
	Trigger      Trigger
	Payload      json.RawMessage
	ScheduledFor *time.Time
	ObservedAt   time.Time
	InitiatedBy  string
}

type Run struct {
	ID           RunID
	TrailID      ID
	InitiatedBy  *xid.ID
	Kind         RunKind
	Trigger      TriggerEvent
	ScheduledFor *time.Time
	Status       RunStatus
	FinishedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ActionRuns   []*ActionRun
}

type ActionRun struct {
	ID             ActionRunID
	RunID          RunID
	ActionID       ActionID
	Action         *Action
	Kind           ActionKind
	Config         json.RawMessage
	Status         ActionRunStatus
	LeaseToken     *string
	LeaseExpiresAt *time.Time
	Output         json.RawMessage
	Target         json.RawMessage
	ErrorText      *string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	NotifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Run            *Run
	Trail          *Trail
}

type ActionSpec struct {
	Kind   ActionKind
	Config json.RawMessage
}

type RobotRunConfig struct {
	Type        ActionKind `json:"type"`
	RobotRef    string     `json:"robot_ref"`
	Instruction string     `json:"instruction"`
}

type RobotInvocationOutput struct {
	Status    RobotInvocationOutputStatus `json:"status"`
	Summary   string                      `json:"summary"`
	Attention *RobotInvocationAttention   `json:"attention,omitempty"`
}

type RobotInvocation struct {
	Type           ActionKind `json:"type"`
	RobotSessionID string     `json:"robot_session_id,omitempty"`
}

type RobotInvocationAttention struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
