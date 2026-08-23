package trail_manager

import (
	"context"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/trail/trail_action_registry"
)

var (
	errInvalidName   = fault.New("trail name must be between 1 and 120 characters", ftag.With(ftag.InvalidArgument))
	errInvalidStatus = fault.New("invalid trail status", ftag.With(ftag.InvalidArgument))
	errArchived      = fault.New("archived trail cannot be changed or run", ftag.With(ftag.InvalidArgument))
)

type Manager struct {
	repository *trail.Repository
	accounts   *account_repo.Repository
	registry   *trail_action_registry.Registry
}

func New(repository *trail.Repository, accounts *account_repo.Repository, registry *trail_action_registry.Registry) *Manager {
	return &Manager{repository: repository, accounts: accounts, registry: registry}
}

func (m *Manager) Create(
	ctx context.Context,
	creatorID xid.ID,
	name string,
	description string,
	status trail.Status,
	trigger trail.Trigger,
	actions []trail.ActionSpec,
) (*trail.Trail, error) {
	validated, next, status, err := m.validate(ctx, creatorID, name, description, status, trigger, actions)
	if err != nil {
		return nil, err
	}

	return m.repository.Create(ctx, creatorID, validated.Name, validated.Description, status, trigger, actions, opt.NewPtr(next))
}

func (m *Manager) Update(
	ctx context.Context,
	id trail.ID,
	name string,
	description string,
	status trail.Status,
	trigger trail.Trigger,
	actions []trail.ActionSpec,
) (*trail.Trail, error) {
	current, err := m.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if current.Status == trail.StatusArchived {
		return nil, errArchived
	}

	validated, next, status, err := m.validate(ctx, xid.ID(current.Creator.ID), name, description, status, trigger, actions)
	if err != nil {
		return nil, err
	}

	return m.repository.Update(ctx, id, validated.Name, validated.Description, status, trigger, actions, opt.NewPtr(next))
}

type validatedWrite struct {
	Name        string
	Description string
}

func (m *Manager) validate(
	ctx context.Context,
	creatorID xid.ID,
	name string,
	description string,
	status trail.Status,
	trigger trail.Trigger,
	actions []trail.ActionSpec,
) (validatedWrite, *time.Time, trail.Status, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 120 {
		return validatedWrite{}, nil, trail.Status{}, errInvalidName
	}

	schedule := trigger.Schedule()
	switch trigger.Type() {
	case trail.TriggerTypeSchedule:
		if schedule == nil {
			return validatedWrite{}, nil, trail.Status{}, trail.ErrInvalidScheduleTrigger
		}

	case trail.TriggerTypeEvent:
		if trigger.Event() == nil {
			return validatedWrite{}, nil, trail.Status{}, trail.ErrInvalidEventTrigger
		}

	default:
		return validatedWrite{}, nil, trail.Status{}, trail.ErrUnsupportedTrigger
	}

	if len(actions) == 0 {
		return validatedWrite{}, nil, trail.Status{}, trail.ErrNoActions
	}

	creator, err := m.accounts.GetRefByID(ctx, account.AccountID(creatorID))
	if err != nil {
		return validatedWrite{}, nil, trail.Status{}, fault.Wrap(err, fctx.With(ctx))
	}

	if err := creator.RejectSuspended(); err != nil {
		return validatedWrite{}, nil, trail.Status{}, fault.Wrap(err, fctx.With(ctx))
	}

	for _, action := range actions {
		if err := m.registry.Validate(ctx, creator, action); err != nil {
			return validatedWrite{}, nil, trail.Status{}, err
		}
	}

	if status == (trail.Status{}) {
		status = trail.StatusActive
	}

	if status != trail.StatusActive && status != trail.StatusPaused && status != trail.StatusArchived {
		return validatedWrite{}, nil, trail.Status{}, errInvalidStatus
	}

	if status == trail.StatusActive {
		if err := creator.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionUseRobots); err != nil {
			return validatedWrite{}, nil, trail.Status{}, fault.Wrap(err, fctx.With(ctx))
		}
	}

	var next *time.Time
	if schedule != nil {
		occurrences := recurrence.Preview(schedule, time.Now().UTC(), 1)
		if len(occurrences) > 0 {
			next = &occurrences[0]
		} else if status == trail.StatusActive {
			status = trail.StatusFinished
		}
	}

	return validatedWrite{Name: name, Description: strings.TrimSpace(description)}, next, status, nil
}

func (m *Manager) Preview(schedule *recurrence.Schedule, after time.Time) []time.Time {
	return recurrence.Preview(schedule, after, 5)
}

func (m *Manager) List(ctx context.Context) ([]*trail.Trail, error) {
	return m.repository.List(ctx)
}

func (m *Manager) Get(ctx context.Context, id trail.ID) (*trail.Trail, error) {
	definition, err := m.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := authoriseManager(ctx); err != nil {
		return nil, err
	}

	return definition, nil
}

func (m *Manager) ListRuns(ctx context.Context, id trail.ID) ([]*trail.Run, error) {
	return m.ListRunsLimited(ctx, id, 100)
}

func (m *Manager) ListRunsLimited(ctx context.Context, id trail.ID, limit int) ([]*trail.Run, error) {
	_, err := m.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := authoriseManager(ctx); err != nil {
		return nil, err
	}

	return m.repository.ListRuns(ctx, id, limit)
}

func (m *Manager) GetRun(ctx context.Context, trailID trail.ID, runID trail.RunID) (*trail.Run, error) {
	run, err := m.repository.GetRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	if run.TrailID != trailID {
		return nil, trail.ErrNotFound
	}

	if err := authoriseManager(ctx); err != nil {
		return nil, err
	}

	return run, nil
}

func (m *Manager) RunNow(ctx context.Context, id trail.ID, initiatedBy xid.ID) (*trail.Run, error) {
	return m.repository.MaterialiseManual(ctx, id, initiatedBy)
}

func authoriseManager(ctx context.Context) error {
	viewer, err := session.GetAccount(ctx)
	if err != nil {
		return err
	}

	return viewer.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionManageTrails)
}

func (m *Manager) Cancel(ctx context.Context, trailID trail.ID, runID trail.RunID, actionRunID trail.ActionRunID) error {
	run, err := m.GetRun(ctx, trailID, runID)
	if err != nil {
		return err
	}

	for _, action := range run.ActionRuns {
		if action.ID == actionRunID {
			return m.registry.Cancel(ctx, action)
		}
	}

	return trail.ErrActionRunNotFound
}
