package robot

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/session_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type SessionID xid.ID

func (id SessionID) String() string {
	return xid.ID(id).String()
}

func NewSessionID(s string) (SessionID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return SessionID{}, err
	}
	return SessionID(id), nil
}

type Session struct {
	session_ref.Ref
	Messages Messages
	State    map[string]any
}

func MapSession(s *ent.RobotSession, messages []*ent.RobotSessionMessage) (*Session, error) {
	user, err := account.MapRef(s.Edges.User)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	mappedMessages, err := dt.MapErr(messages, MapMessage)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Session{
		Ref: session_ref.Ref{
			ID:        session_ref.ID(s.ID),
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			Human:     *user,
		},
		Messages: mappedMessages,
		State:    s.State,
	}, nil
}

func MapSessionRef(s *ent.RobotSession) (*session_ref.Ref, error) {
	user, err := account.MapRef(s.Edges.User)
	if err != nil {
		return nil, err
	}

	return &session_ref.Ref{
		ID:        session_ref.ID(s.ID),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		Name:      s.Name,
		Human:     *user,
	}, nil
}
