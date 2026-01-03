package robot

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
)

type SessionID xid.ID

func (id SessionID) String() string {
	return xid.ID(id).String()
}

type Session struct {
	ID        SessionID
	CreatedAt time.Time
	UpdatedAt time.Time

	RobotID opt.Optional[robot_ref.ID]
	UserID  account.AccountID
	State   map[string]any
}

type MessageID xid.ID

func (id MessageID) String() string {
	return xid.ID(id).String()
}

type Message struct {
	ID        MessageID
	CreatedAt time.Time

	SessionID    SessionID
	InvocationID string
	AuthorID     opt.Optional[account.AccountID]
	EventData    map[string]any
}

type Messages []*Message
