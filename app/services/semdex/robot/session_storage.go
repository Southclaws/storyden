package robot

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
)

const currentRobotIDStateKey = "current_robot_id"

type sessionStorage struct {
	robotQuerier     *robot_querier.Querier
	robotSessionRepo *robot_session.Repository
}

func NewSessionStorage(robotQuerier *robot_querier.Querier, robotSessionRepo *robot_session.Repository) adksession.Service {
	return &sessionStorage{
		robotQuerier:     robotQuerier,
		robotSessionRepo: robotSessionRepo,
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

	name := getDefaultSessionName()

	sess, err := s.robotSessionRepo.Create(ctx, robot.SessionID(sessionID), name, account.AccountID(userID), req.State)
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

	sess, err := s.robotSessionRepo.GetWithMessageFilters(
		ctx,
		robot.SessionID(sessionID),
		opt.NewEmpty[account.AccountID](),
		req.NumRecentEvents,
		req.After,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	events := mapToADKEventsFromMessages(sess.Messages)

	return &adksession.GetResponse{
		Session: s.mapToADKSessionForUser(sess, events, req.AppName, userID),
	}, nil
}

func (s *sessionStorage) List(ctx context.Context, req *adksession.ListRequest) (*adksession.ListResponse, error) {
	if req.UserID == "" {
		return &adksession.ListResponse{
			Sessions: []adksession.Session{},
		}, nil
	}

	userID, err := xid.FromString(req.UserID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sessions, err := s.robotSessionRepo.ListAll(ctx, opt.New(account.AccountID(userID)))
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

	return s.robotSessionRepo.Delete(ctx, robot.SessionID(sessionID))
}

func (s *sessionStorage) AppendEvent(ctx context.Context, sess adksession.Session, event *adksession.Event) error {
	if event.Partial {
		return nil
	}
	if isClientSidePendingEvent(event) {
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

	var accountIDOpt opt.Optional[account.AccountID]
	if event.Author == "user" {
		accountID, err := xid.FromString(sess.UserID())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		accountIDOpt = opt.New(account.AccountID(accountID))
	}

	actorOpt := opt.NewEmpty[robot.Actor]()
	if event.Author != "user" {
		actor, err := s.actorForEvent(ctx, sess, event)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		actorOpt = actor
	}

	eventData, err := structToMap(event)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.robotSessionRepo.AppendMessage(
		ctx,
		robot.SessionID(sessionID),
		event.InvocationID,
		accountIDOpt,
		actorOpt,
		eventData,
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ourSession.events = append(ourSession.events.(adkEvents), event)
	ourSession.lastUpdateTime = event.Timestamp

	return nil
}

func (s *sessionStorage) actorForEvent(ctx context.Context, sess adksession.Session, event *adksession.Event) (opt.Optional[robot.Actor], error) {
	if actor, ok, err := actorFromSessionState(sess.State()); err != nil {
		return opt.NewEmpty[robot.Actor](), err
	} else if ok {
		return opt.New(actor), nil
	}

	if event.Author == "storyden" {
		return opt.New(robot.NewBuiltinActor(agent_registry.RobotBuilderID)), nil
	}

	dbRobot, err := s.robotQuerier.GetByName(ctx, event.Author)
	if err == nil {
		return opt.New(robot.NewDatabaseActor(xid.ID(dbRobot.ID))), nil
	}
	if ftag.Get(err) != ftag.NotFound {
		return opt.NewEmpty[robot.Actor](), err
	}

	return opt.New(robot.NewBuiltinActor(event.Author)), nil
}

func actorFromSessionState(state adksession.State) (robot.Actor, bool, error) {
	if state == nil {
		return robot.Actor{}, false, nil
	}

	raw, err := state.Get(currentRobotIDStateKey)
	if err != nil {
		return robot.Actor{}, false, nil
	}

	ref, ok := raw.(string)
	if !ok || strings.TrimSpace(ref) == "" {
		return robot.Actor{}, false, nil
	}

	actor, err := robot.NewActorFromRobotRef(ref)
	if err != nil {
		return robot.Actor{}, false, err
	}

	return actor, true, nil
}

func CurrentRobotIDFromState(state map[string]any) opt.Optional[string] {
	raw, ok := state[currentRobotIDStateKey]
	if !ok {
		return opt.NewEmpty[string]()
	}

	ref, ok := raw.(string)
	if !ok || strings.TrimSpace(ref) == "" {
		return opt.NewEmpty[string]()
	}

	return opt.New(ref)
}

func isClientSidePendingEvent(event *adksession.Event) bool {
	if event == nil || event.LLMResponse.Content == nil {
		return false
	}

	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || part.FunctionResponse == nil {
			continue
		}
		pending, _ := part.FunctionResponse.Response["_client_side_pending"].(bool)
		if pending {
			return true
		}
	}

	return false
}

func (s *sessionStorage) mapToADKSession(sess *robot.Session, events adksession.Events, appName string) adksession.Session {
	return s.mapToADKSessionForUser(sess, events, appName, xid.ID(sess.Human.ID))
}

func (s *sessionStorage) mapToADKSessionForUser(sess *robot.Session, events adksession.Events, appName string, userID xid.ID) adksession.Session {
	if events == nil {
		events = adkEvents{}
	}

	state := adkState(sess.State)
	if state == nil {
		state = make(adkState)
	}

	return &adkSession{
		id:             xid.ID(sess.ID).String(),
		appName:        appName,
		userID:         userID.String(),
		state:          state,
		events:         events,
		lastUpdateTime: sess.UpdatedAt,
	}
}

func mapToADKEventsFromMessages(messages []*robot.Message) adksession.Events {
	events := dt.Map(messages, func(m *robot.Message) *adksession.Event {
		return &m.Event
	})

	return adkEvents(events)
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

func getDefaultSessionName() string {
	return "Untitled (" + time.Now().Format("January 2, 2006 at 3:04 PM") + ")"
}
