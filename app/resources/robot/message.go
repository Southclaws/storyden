package robot

import (
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type MessageID xid.ID

func (id MessageID) String() string {
	return xid.ID(id).String()
}

// Message roughly maps to OUR database row.
type Message struct {
	ID        MessageID
	CreatedAt time.Time

	SessionID    SessionID
	InvocationID string
	Robot        opt.Optional[*Robot]
	Author       opt.Optional[*account.Account]

	Event session.Event
}

type Messages []*Message

func MapMessage(m *ent.RobotSessionMessage) (*Message, error) {
	var robotOpt opt.Optional[*Robot]
	if m.Edges.Robot != nil {
		r, err := Map(m.Edges.Robot)
		if err != nil {
			return nil, err
		}
		robotOpt = opt.New(r)
	}

	var authorOpt opt.Optional[*account.Account]
	if m.Edges.Author != nil {
		a, err := account.MapRef(m.Edges.Author)
		if err != nil {
			return nil, err
		}
		authorOpt = opt.New(a)
	}

	evt, err := mapToADKEventFromMessage(m.EventData)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Message{
		ID:           MessageID(m.ID),
		CreatedAt:    m.CreatedAt,
		SessionID:    SessionID(m.SessionID),
		InvocationID: m.InvocationID,
		Robot:        robotOpt,
		Author:       authorOpt,
		Event:        *evt,
	}, nil
}

func mapToADKEventFromMessage(data map[string]any) (*session.Event, error) {
	var event session.Event
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
