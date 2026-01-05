package robot

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/session"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/robot"
	"github.com/Southclaws/storyden/internal/ent/robotsession"
	"github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
)

type sessionStorage struct {
	db *ent.Client
}

func NewSessionStorage(db *ent.Client) adksession.Service {
	return &sessionStorage{
		db: db,
	}
}

// adkSession implements the session.Session interface
type adkSession struct {
	id             string
	appName        string
	userID         string
	state          adksession.State
	events         adksession.Events
	lastUpdateTime time.Time
}

func (s *adkSession) ID() string                { return s.id }
func (s *adkSession) AppName() string           { return s.appName }
func (s *adkSession) UserID() string            { return s.userID }
func (s *adkSession) State() adksession.State   { return s.state }
func (s *adkSession) Events() adksession.Events { return s.events }
func (s *adkSession) LastUpdateTime() time.Time { return s.lastUpdateTime }

// adkState implements the session.State interface
type adkState map[string]any

func (s adkState) Get(key string) (any, error) {
	val, ok := s[key]
	if !ok {
		return nil, adksession.ErrStateKeyNotExist
	}
	return val, nil
}

func (s adkState) Set(key string, value any) error {
	s[key] = value
	return nil
}

func (s adkState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range s {
			if !yield(k, v) {
				return
			}
		}
	}
}

// adkEvents implements the session.Events interface
type adkEvents []*adksession.Event

func (e adkEvents) All() iter.Seq[*adksession.Event] {
	return func(yield func(*adksession.Event) bool) {
		for _, event := range e {
			if !yield(event) {
				return
			}
		}
	}
}

func (e adkEvents) Len() int {
	return len(e)
}

func (e adkEvents) At(i int) *adksession.Event {
	if i < 0 || i >= len(e) {
		return nil
	}
	return e[i]
}

func (s *sessionStorage) Create(ctx context.Context, req *adksession.CreateRequest) (*adksession.CreateResponse, error) {
	// Parse IDs
	userID, err := xid.FromString(req.UserID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Parse or generate session ID
	var sessionID xid.ID
	if req.SessionID != "" {
		sessionID, err = xid.FromString(req.SessionID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		sessionID = xid.New()
	}

	sess, err := s.db.RobotSession.Create().
		SetID(sessionID).
		SetAccountID(userID).
		SetState(req.State).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &adksession.CreateResponse{
		Session: s.mapToADKSession(sess, nil, req.AppName),
	}, nil
}

func (s *sessionStorage) Get(ctx context.Context, req *adksession.GetRequest) (*adksession.GetResponse, error) {
	sessionID, err := xid.FromString(req.SessionID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	userID, err := xid.FromString(req.UserID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := s.db.RobotSession.Query().
		Where(
			robotsession.IDEQ(sessionID),
			robotsession.AccountIDEQ(userID),
		).
		WithMessages(func(q *ent.RobotSessionMessageQuery) {
			q.Order(ent.Asc(robotsessionmessage.FieldCreatedAt))

			// Apply filters
			if req.NumRecentEvents > 0 {
				q.Limit(req.NumRecentEvents)
			}
			if !req.After.IsZero() {
				q.Where(robotsessionmessage.CreatedAtGTE(req.After))
			}
		})

	sess, err := query.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	messages, err := sess.Edges.MessagesOrErr()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	events, err := mapToADKEvents(messages)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &adksession.GetResponse{
		Session: s.mapToADKSession(sess, events, req.AppName),
	}, nil
}

func (s *sessionStorage) List(ctx context.Context, req *adksession.ListRequest) (*adksession.ListResponse, error) {
	query := s.db.RobotSession.Query()

	if req.UserID != "" {
		userID, err := xid.FromString(req.UserID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		query = query.Where(robotsession.AccountIDEQ(userID))
	}

	sessions, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := make([]adksession.Session, len(sessions))
	for i, sess := range sessions {
		result[i] = s.mapToADKSession(sess, nil, req.AppName)
	}

	return &adksession.ListResponse{
		Sessions: result,
	}, nil
}

func (s *sessionStorage) Delete(ctx context.Context, req *adksession.DeleteRequest) error {
	sessionID, err := xid.FromString(req.SessionID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.db.RobotSession.DeleteOneID(sessionID).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *sessionStorage) AppendEvent(ctx context.Context, sess adksession.Session, event *adksession.Event) error {
	if event.Partial {
		return nil
	}

	sessionID, err := xid.FromString(sess.ID())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ourSession, ok := sess.(*adkSession)
	if !ok {
		return fault.Wrap(fmt.Errorf("unexpected session type %T", sess), fctx.With(ctx))
	}

	create := s.db.RobotSessionMessage.Create().
		SetSessionID(sessionID).
		SetInvocationID(event.InvocationID)

	// Set author_id for user messages
	if event.Author == "user" {
		accountID, err := xid.FromString(sess.UserID())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		create.SetAccountID(accountID)
	}

	// Set robot_id for agent messages
	// event.Author contains the agent name from llmagent.Config.Name
	// Look up robot by name to get the robot_id
	if event.Author != "user" && event.Author != "" {
		// Query robot by name
		rb, err := s.db.Robot.Query().
			Where(robot.NameEQ(event.Author)).
			Only(ctx)

		// If robot found, set robot_id
		// If not found (e.g., "storyden" default agent), leave robot_id as null
		if err == nil {
			create.SetRobotID(rb.ID)
		} else if !ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	eventData, err := structToMap(event)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	create.SetEventData(eventData)

	_, err = create.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ourSession.events = append(ourSession.events.(adkEvents), event)
	ourSession.lastUpdateTime = event.Timestamp

	return nil
}

func (s *sessionStorage) mapToADKSession(sess *ent.RobotSession, events adksession.Events, appName string) adksession.Session {
	if events == nil {
		events = adkEvents{}
	}

	state := adkState(sess.State)
	if state == nil {
		state = make(adkState)
	}

	return &adkSession{
		id:             sess.ID.String(),
		appName:        appName,
		userID:         sess.AccountID.String(),
		state:          state,
		events:         events,
		lastUpdateTime: sess.UpdatedAt,
	}
}

func mapToADKEvents(messages []*ent.RobotSessionMessage) (adksession.Events, error) {
	events, err := dt.MapErr(messages, mapToADKEvent)
	if err != nil {
		return nil, err
	}
	return adkEvents(events), nil
}

func mapToADKEvent(msg *ent.RobotSessionMessage) (*adksession.Event, error) {
	var event adksession.Event
	data, err := json.Marshal(msg.EventData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func structToMap(v any) (map[string]any, error) {
	if v == nil {
		return nil, nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
