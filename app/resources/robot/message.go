package robot

import (
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/session"

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
	Actor        opt.Optional[Actor]
	Author       opt.Optional[*account.Account]

	Event session.Event
}

type Messages []*Message

type MessageCursorParams struct {
	Before opt.Optional[MessageID]
	Size   int
}

func NewMessageCursorParams(before opt.Optional[MessageID], size int) MessageCursorParams {
	if size <= 0 {
		size = 50
	}

	return MessageCursorParams{
		Before: before,
		Size:   size,
	}
}

func (p MessageCursorParams) Limit() int {
	return p.Size
}

type MessageCursorResult struct {
	Size       int
	Results    int
	NextBefore opt.Optional[MessageID]
	Items      Messages
}

func MapMessage(m *ent.RobotSessionMessage) (*Message, error) {
	var robotOpt opt.Optional[*Robot]
	if m.Edges.Robot != nil {
		r, err := Map(m.Edges.Robot)
		if err != nil {
			return nil, err
		}
		robotOpt = opt.New(r)
	}

	actorOpt, err := MapMessageActor(m, robotOpt)
	if err != nil {
		return nil, err
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
		Actor:        actorOpt,
		Author:       authorOpt,
		Event:        *evt,
	}, nil
}

func MapMessageActor(m *ent.RobotSessionMessage, robotOpt opt.Optional[*Robot]) (opt.Optional[Actor], error) {
	if m.RobotID != nil && m.BuiltinRobot != nil {
		return opt.NewEmpty[Actor](), fault.New("robot session message has both database and built-in robot actors")
	}

	if m.RobotID != nil {
		return opt.New(NewDatabaseActor(*m.RobotID)), nil
	}

	if robot, ok := robotOpt.Get(); ok {
		return opt.New(NewDatabaseActor(xid.ID(robot.ID))), nil
	}

	if m.BuiltinRobot != nil {
		return opt.New(NewBuiltinActor(*m.BuiltinRobot)), nil
	}

	return opt.NewEmpty[Actor](), nil
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
