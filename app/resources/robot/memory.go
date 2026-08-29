package robot

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

//go:generate go run github.com/Southclaws/enumerator

type MemoryID xid.ID

func (id MemoryID) String() string { return xid.ID(id).String() }

func NewMemoryID(value string) (MemoryID, error) {
	id, err := xid.FromString(value)
	return MemoryID(id), err
}

type memoryStateEnum string

const (
	memoryStateActive     memoryStateEnum = "active"
	memoryStateSuperseded memoryStateEnum = "superseded"
	memoryStateArchived   memoryStateEnum = "archived"
)

type Memory struct {
	ID             MemoryID
	RobotRef       string
	ParentID       opt.Optional[MemoryID]
	Content        string
	Fact           opt.Optional[MemoryFact]
	State          MemoryState
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastAccessedAt time.Time
	AccessCount    uint64
}

type MemoryFact struct {
	Subject   string
	Predicate string
	Object    string
}

func MapMemory(in *ent.RobotMemory) *Memory {
	fact := opt.NewEmpty[MemoryFact]()
	if in.Subject != nil && in.Predicate != nil && in.Object != nil {
		fact = opt.New(MemoryFact{Subject: *in.Subject, Predicate: *in.Predicate, Object: *in.Object})
	}
	return &Memory{
		ID:       MemoryID(in.ID),
		RobotRef: in.RobotRef,
		ParentID: opt.Map(opt.NewPtr(in.ParentID), func(id xid.ID) MemoryID {
			return MemoryID(id)
		}),
		Content:        in.Content,
		Fact:           fact,
		State:          MemoryState{memoryStateEnum(in.State)},
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		LastAccessedAt: in.LastAccessedAt,
		AccessCount:    in.AccessCount,
	}
}
