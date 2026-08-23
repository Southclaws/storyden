package trail

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	enttrail "github.com/Southclaws/storyden/internal/ent/trail"
	entaction "github.com/Southclaws/storyden/internal/ent/trailaction"
	entactionrun "github.com/Southclaws/storyden/internal/ent/trailactionrun"
	entrun "github.com/Southclaws/storyden/internal/ent/trailrun"
	entscheduler "github.com/Southclaws/storyden/internal/ent/trailschedulerlease"
)

var (
	ErrNotFound           = fault.New("trail not found", ftag.With(ftag.NotFound))
	ErrNoActions          = fault.New("trail requires at least one action", ftag.With(ftag.InvalidArgument))
	ErrActionRunNotFound  = fault.New("trail action run not found", ftag.With(ftag.NotFound))
	ErrUnsupportedTrigger = fault.New("unsupported trail trigger", ftag.With(ftag.InvalidArgument))
)

type Repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	creatorID xid.ID,
	name string,
	description string,
	status Status,
	trigger Trigger,
	actions []ActionSpec,
	nextOccurrenceAt opt.Optional[time.Time],
) (*Trail, error) {
	if len(actions) == 0 {
		return nil, ErrNoActions
	}
	triggerType, triggerConfig, err := encodeTrigger(trigger)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	id := xid.New()
	err = ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		create := tx.Trail.Create().
			SetID(id).
			SetAccountID(creatorID).
			SetName(name).
			SetDescription(description).
			SetStatus(enttrail.Status(status.String())).
			SetTriggerType(triggerType.String()).
			SetTriggerConfig(triggerConfig).
			SetNillableNextOccurrenceAt(nextOccurrenceAt.Ptr())
		if _, err := create.Save(ctx); err != nil {
			return err
		}
		for position, action := range actions {
			if _, err := tx.TrailAction.Create().
				SetTrailID(id).
				SetKind(action.Kind.String()).
				SetPosition(position).
				SetConfig(action.Config).
				Save(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return r.Get(ctx, ID(id))
}

func (r *Repository) Update(
	ctx context.Context,
	id ID,
	name string,
	description string,
	status Status,
	trigger Trigger,
	actions []ActionSpec,
	nextOccurrenceAt opt.Optional[time.Time],
) (*Trail, error) {
	if len(actions) == 0 {
		return nil, ErrNoActions
	}
	triggerType, triggerConfig, err := encodeTrigger(trigger)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	now := time.Now().UTC()
	err = ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		update := tx.Trail.Update().
			Where(enttrail.IDEQ(xid.ID(id))).
			SetName(name).
			SetDescription(description).
			SetStatus(enttrail.Status(status.String())).
			SetTriggerType(triggerType.String()).
			SetTriggerConfig(triggerConfig)
		if next := nextOccurrenceAt.Ptr(); next != nil {
			update.SetNextOccurrenceAt(*next)
		} else {
			update.ClearNextOccurrenceAt()
		}
		updated, err := update.Save(ctx)
		if err != nil {
			return err
		}
		if updated == 0 {
			return ErrNotFound
		}
		if _, err := tx.TrailAction.Update().
			Where(entaction.TrailIDEQ(xid.ID(id)), entaction.ArchivedAtIsNil()).
			SetArchivedAt(now).
			Save(ctx); err != nil {
			return err
		}
		for position, action := range actions {
			if _, err := tx.TrailAction.Create().
				SetTrailID(xid.ID(id)).
				SetKind(action.Kind.String()).
				SetPosition(position).
				SetConfig(action.Config).
				Save(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return r.Get(ctx, id)
}

func activeActions(query *ent.TrailActionQuery) {
	query.Where(entaction.ArchivedAtIsNil()).Order(entaction.ByPosition(sql.OrderAsc()))
}

func withCreator(query *ent.AccountQuery) {
	query.WithAccountRoles(func(roles *ent.AccountRolesQuery) {
		roles.WithRole()
	})
}

func (r *Repository) Get(ctx context.Context, id ID) (*Trail, error) {
	row, err := r.db.Trail.Query().
		Where(enttrail.IDEQ(xid.ID(id))).
		WithCreator(withCreator).
		WithActions(activeActions).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return mapTrail(row)
}

func (r *Repository) List(ctx context.Context) ([]*Trail, error) {
	rows, err := r.db.Trail.Query().
		WithCreator(withCreator).
		WithActions(activeActions).
		Order(enttrail.ByUpdatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	result := make([]*Trail, len(rows))
	for i, row := range rows {
		result[i], err = mapTrail(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}
	return result, nil
}

type EventTriggerDefinition struct {
	ID     ID
	Events []string
}

func (r *Repository) ActiveEventTriggers(ctx context.Context) ([]EventTriggerDefinition, error) {
	rows, err := r.db.Trail.Query().
		Where(
			enttrail.StatusEQ(enttrail.StatusActive),
			enttrail.TriggerTypeEQ(TriggerTypeEvent.String()),
		).
		Select(enttrail.FieldID, enttrail.FieldTriggerType, enttrail.FieldTriggerConfig).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := make([]EventTriggerDefinition, len(rows))
	for i, row := range rows {
		trigger, err := decodeTrigger(row.TriggerType, row.TriggerConfig)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		event := trigger.Event()
		if event == nil {
			return nil, fault.Wrap(ErrInvalidEventTrigger, fctx.With(ctx))
		}
		result[i] = EventTriggerDefinition{ID: ID(row.ID), Events: event.Events}
	}

	return result, nil
}

func (r *Repository) Due(ctx context.Context, now time.Time, limit int) ([]*Trail, error) {
	rows, err := r.db.Trail.Query().
		Where(
			enttrail.StatusEQ(enttrail.StatusActive),
			enttrail.TriggerTypeEQ(TriggerTypeSchedule.String()),
			enttrail.NextOccurrenceAtNotNil(),
			enttrail.NextOccurrenceAtLTE(now.UTC()),
		).
		WithCreator(withCreator).
		WithActions(activeActions).
		Order(enttrail.ByNextOccurrenceAt(sql.OrderAsc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	result := make([]*Trail, len(rows))
	for i, row := range rows {
		result[i], err = mapTrail(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}
	return result, nil
}

func (r *Repository) MaterialiseScheduled(ctx context.Context, id ID, expectedCurrent, recordedFor time.Time, next *time.Time, finished bool, skipped bool) (*Run, bool, error) {
	var runID xid.ID
	created := false
	err := ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		trailRow, err := tx.Trail.Query().
			Where(
				enttrail.IDEQ(xid.ID(id)),
				enttrail.StatusEQ(enttrail.StatusActive),
				enttrail.NextOccurrenceAtEQ(expectedCurrent.UTC()),
			).
			WithActions(activeActions).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil
			}
			return err
		}
		runID = xid.New()
		definition, err := decodeTrigger(trailRow.TriggerType, trailRow.TriggerConfig)
		if err != nil {
			return err
		}
		event := TriggerEvent{
			TrailID:      id.String(),
			TrailRunID:   runID.String(),
			Kind:         RunKindScheduled,
			Trigger:      definition,
			ScheduledFor: ptrTime(recordedFor.UTC()),
			ObservedAt:   time.Now().UTC(),
		}
		payload, err := encodeTriggerEvent(event)
		if err != nil {
			return err
		}
		status := entrun.StatusQueued
		var finishedAt *time.Time
		if skipped {
			status = entrun.StatusSkipped
			finishedAt = ptrTime(time.Now().UTC())
		}
		if _, err := tx.TrailRun.Create().
			SetID(runID).
			SetTrailID(xid.ID(id)).
			SetKind(entrun.KindScheduled).
			SetTriggerPayload(payload).
			SetScheduledFor(recordedFor.UTC()).
			SetStatus(status).
			SetNillableFinishedAt(finishedAt).
			Save(ctx); err != nil {
			if ent.IsConstraintError(err) {
				return nil
			}
			return err
		}
		if !skipped {
			for _, action := range trailRow.Edges.Actions {
				if _, err := tx.TrailActionRun.Create().
					SetRunID(runID).
					SetActionID(action.ID).
					SetKind(action.Kind).
					SetConfig(action.Config).
					SetStatus(entactionrun.StatusQueued).
					Save(ctx); err != nil {
					return err
				}
			}
		}
		update := tx.Trail.UpdateOneID(xid.ID(id)).
			SetLastOccurrenceAt(recordedFor.UTC()).
			SetNillableNextOccurrenceAt(next)
		if next == nil && finished {
			update.SetStatus(enttrail.StatusFinished).ClearNextOccurrenceAt()
		}
		if _, err := update.Save(ctx); err != nil {
			return err
		}
		created = true
		return nil
	})
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if !created {
		return nil, false, nil
	}
	run, err := r.GetRun(ctx, RunID(runID))
	return run, true, err
}

func (r *Repository) MaterialiseEvent(ctx context.Context, id ID, eventName string, sourcePayload json.RawMessage, observedAt time.Time) (RunID, bool, error) {
	var runID xid.ID
	created := false
	err := ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		trailRow, err := tx.Trail.Query().
			Where(
				enttrail.IDEQ(xid.ID(id)),
				enttrail.StatusEQ(enttrail.StatusActive),
				enttrail.TriggerTypeEQ(TriggerTypeEvent.String()),
			).
			WithActions(activeActions).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil
			}
			return err
		}

		definition, err := decodeTrigger(trailRow.TriggerType, trailRow.TriggerConfig)
		if err != nil {
			return err
		}
		eventTrigger := definition.Event()
		if eventTrigger == nil || !slices.Contains(eventTrigger.Events, eventName) {
			return nil
		}
		sourcePayload, err = materialiseEventPayload(eventName, sourcePayload)
		if err != nil {
			return err
		}

		runID = xid.New()
		event := TriggerEvent{
			TrailID:    id.String(),
			TrailRunID: runID.String(),
			Kind:       RunKindEvent,
			EventName:  eventName,
			Trigger:    definition,
			Payload:    sourcePayload,
			ObservedAt: observedAt.UTC(),
		}
		payload, err := encodeTriggerEvent(event)
		if err != nil {
			return err
		}

		// Event occurrences are automatic runs. The existing database kind keeps
		// that distinction from manual runs while the immutable trigger snapshot
		// carries the more specific event kind.
		if _, err := tx.TrailRun.Create().
			SetID(runID).
			SetTrailID(xid.ID(id)).
			SetKind(entrun.KindScheduled).
			SetTriggerPayload(payload).
			SetStatus(entrun.StatusQueued).
			Save(ctx); err != nil {
			return err
		}
		for _, action := range trailRow.Edges.Actions {
			if _, err := tx.TrailActionRun.Create().
				SetRunID(runID).
				SetActionID(action.ID).
				SetKind(action.Kind).
				SetConfig(action.Config).
				SetStatus(entactionrun.StatusQueued).
				Save(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.Trail.UpdateOneID(xid.ID(id)).SetLastOccurrenceAt(observedAt.UTC()).Save(ctx); err != nil {
			return err
		}

		created = true
		return nil
	})
	if err != nil {
		return RunID{}, false, fault.Wrap(err, fctx.With(ctx))
	}
	if !created {
		return RunID{}, false, nil
	}

	return RunID(runID), true, nil
}

func (r *Repository) MaterialiseManual(ctx context.Context, id ID, initiatedBy xid.ID) (*Run, error) {
	trailRow, err := r.db.Trail.Query().Where(
		enttrail.IDEQ(xid.ID(id)),
		enttrail.StatusNEQ(enttrail.StatusArchived),
	).WithActions(activeActions).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	runID := xid.New()
	definition, err := decodeTrigger(trailRow.TriggerType, trailRow.TriggerConfig)
	if err != nil {
		return nil, err
	}
	event := TriggerEvent{
		TrailID:     id.String(),
		TrailRunID:  runID.String(),
		Kind:        RunKindManual,
		Trigger:     definition,
		ObservedAt:  time.Now().UTC(),
		InitiatedBy: initiatedBy.String(),
	}
	payload, err := encodeTriggerEvent(event)
	if err != nil {
		return nil, err
	}
	err = ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		if _, err := tx.TrailRun.Create().
			SetID(runID).
			SetTrailID(xid.ID(id)).
			SetInitiatedByID(initiatedBy).
			SetKind(entrun.KindManual).
			SetTriggerPayload(payload).
			SetStatus(entrun.StatusQueued).
			Save(ctx); err != nil {
			return err
		}
		for _, action := range trailRow.Edges.Actions {
			if _, err := tx.TrailActionRun.Create().
				SetRunID(runID).
				SetActionID(action.ID).
				SetKind(action.Kind).
				SetConfig(action.Config).
				SetStatus(entactionrun.StatusQueued).
				Save(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return r.GetRun(ctx, RunID(runID))
}

func withActionRuns(query *ent.TrailActionRunQuery) {
	query.WithAction().Order(entactionrun.ByCreatedAt(sql.OrderAsc()))
}

func (r *Repository) GetRun(ctx context.Context, id RunID) (*Run, error) {
	row, err := r.db.TrailRun.Query().
		Where(entrun.IDEQ(xid.ID(id))).
		WithActionRuns(withActionRuns).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	mapped, err := mapRun(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return mapped, nil
}

// RunNumber returns the one-based position of a run in its Trail history.
// Trail runs are append-only, so the chronological position remains stable.
func (r *Repository) RunNumber(ctx context.Context, id RunID) (int, error) {
	row, err := r.db.TrailRun.Query().
		Where(entrun.IDEQ(xid.ID(id))).
		Select(entrun.FieldID, entrun.FieldTrailID, entrun.FieldCreatedAt).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, ErrNotFound
		}
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	number, err := r.db.TrailRun.Query().Where(
		entrun.TrailIDEQ(row.TrailID),
		entrun.Or(
			entrun.CreatedAtLT(row.CreatedAt),
			entrun.And(entrun.CreatedAtEQ(row.CreatedAt), entrun.IDLTE(row.ID)),
		),
	).Count(ctx)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}
	return number, nil
}

func (r *Repository) ListRuns(ctx context.Context, trailID ID, limit int) ([]*Run, error) {
	rows, err := r.db.TrailRun.Query().
		Where(entrun.TrailIDEQ(xid.ID(trailID))).
		WithActionRuns(withActionRuns).
		Order(entrun.ByCreatedAt(sql.OrderDesc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	result := make([]*Run, len(rows))
	for i, row := range rows {
		result[i], err = mapRun(row)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}
	return result, nil
}

type SchedulerLease struct {
	Token     string
	ExpiresAt time.Time
}

func (r *Repository) AcquireSchedulerLease(ctx context.Context, now time.Time, duration time.Duration) (*SchedulerLease, bool, error) {
	const singletonID = 1
	if err := r.db.TrailSchedulerLease.Create().
		SetID(singletonID).
		OnConflictColumns(entscheduler.FieldID).
		Ignore().
		Exec(ctx); err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	row, err := r.db.TrailSchedulerLease.Get(ctx, singletonID)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if row.LeaseToken != nil && row.LeaseExpiresAt != nil && row.LeaseExpiresAt.After(now) {
		return nil, false, nil
	}
	token := uuid.NewString()
	expires := now.Add(duration)
	updated, err := r.db.TrailSchedulerLease.Update().
		Where(
			entscheduler.IDEQ(singletonID),
			entscheduler.Or(entscheduler.LeaseTokenIsNil(), entscheduler.LeaseExpiresAtIsNil(), entscheduler.LeaseExpiresAtLTE(now)),
		).
		SetLeaseToken(token).
		SetLeaseExpiresAt(expires).
		Save(ctx)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if updated == 0 {
		return nil, false, nil
	}
	return &SchedulerLease{Token: token, ExpiresAt: expires}, true, nil
}

func (r *Repository) RenewSchedulerLease(ctx context.Context, token string, expires time.Time) (bool, error) {
	const singletonID = 1
	updated, err := r.db.TrailSchedulerLease.Update().
		Where(entscheduler.IDEQ(singletonID), entscheduler.LeaseTokenEQ(token), entscheduler.LeaseExpiresAtGT(time.Now().UTC())).
		SetLeaseExpiresAt(expires.UTC()).
		Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return updated == 1, nil
}

func (r *Repository) ClaimActionRun(ctx context.Context, now time.Time, leaseDuration time.Duration) (*ActionRun, bool, error) {
	return r.claimActionRun(ctx, nil, now, leaseDuration)
}

func (r *Repository) ClaimActionRunByID(ctx context.Context, id ActionRunID, now time.Time, leaseDuration time.Duration) (*ActionRun, bool, error) {
	return r.claimActionRun(ctx, &id, now, leaseDuration)
}

func (r *Repository) claimActionRun(ctx context.Context, id *ActionRunID, now time.Time, leaseDuration time.Duration) (*ActionRun, bool, error) {
	query := r.db.TrailActionRun.Query().
		Where(entactionrun.Or(
			entactionrun.StatusEQ(entactionrun.StatusQueued),
			entactionrun.And(
				entactionrun.StatusEQ(entactionrun.StatusRunning),
				entactionrun.Or(entactionrun.LeaseExpiresAtIsNil(), entactionrun.LeaseExpiresAtLTE(now)),
			),
		)).
		Order(entactionrun.ByCreatedAt(sql.OrderAsc()))
	if id != nil {
		query.Where(entactionrun.IDEQ(xid.ID(*id)))
	}

	row, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	token := uuid.NewString()
	claimable := entactionrun.StatusEQ(row.Status)
	if row.Status == entactionrun.StatusRunning {
		claimable = entactionrun.And(
			entactionrun.StatusEQ(entactionrun.StatusRunning),
			entactionrun.Or(entactionrun.LeaseExpiresAtIsNil(), entactionrun.LeaseExpiresAtLTE(now)),
		)
	}
	update := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(row.ID),
		claimable,
	).SetStatus(entactionrun.StatusRunning).
		SetLeaseToken(token).
		SetLeaseExpiresAt(now.Add(leaseDuration))
	if row.StartedAt == nil {
		update.SetStartedAt(now)
	}
	updated, err := update.Save(ctx)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if updated == 0 {
		return nil, false, nil
	}
	claimed, err := r.actionRunWithContext(ctx, row.ID)
	if err != nil {
		return nil, false, err
	}
	claimed.LeaseToken = &token
	expires := now.Add(leaseDuration)
	claimed.LeaseExpiresAt = &expires
	claimed.Status = ActionRunStatusRunning
	return claimed, true, nil
}

func (r *Repository) actionRunWithContext(ctx context.Context, id xid.ID) (*ActionRun, error) {
	row, err := r.db.TrailActionRun.Query().
		Where(entactionrun.IDEQ(id)).
		WithAction().
		WithRun(func(query *ent.TrailRunQuery) {
			query.WithTrail(func(tq *ent.TrailQuery) {
				tq.WithCreator(withCreator).WithActions(activeActions)
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrActionRunNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	mapped, err := mapActionRun(row)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if row.Edges.Run != nil {
		mapped.Run, err = mapRun(row.Edges.Run)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		if row.Edges.Run.Edges.Trail != nil {
			mapped.Trail, err = mapTrail(row.Edges.Run.Edges.Trail)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}
	return mapped, nil
}

func (r *Repository) GetActionRun(ctx context.Context, id ActionRunID) (*ActionRun, error) {
	return r.actionRunWithContext(ctx, xid.ID(id))
}

func (r *Repository) SetActionRunTarget(ctx context.Context, id ActionRunID, token string, target json.RawMessage) error {
	updated, err := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(xid.ID(id)), entactionrun.LeaseTokenEQ(token), entactionrun.StatusEQ(entactionrun.StatusRunning),
	).SetTarget(target).Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if updated != 1 {
		return ErrActionRunNotFound
	}
	return nil
}

func (r *Repository) RenewActionLease(ctx context.Context, id ActionRunID, token string, expires time.Time) (bool, error) {
	updated, err := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(xid.ID(id)), entactionrun.LeaseTokenEQ(token), entactionrun.StatusEQ(entactionrun.StatusRunning),
	).SetLeaseExpiresAt(expires).Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return updated == 1, nil
}

func (r *Repository) FinishActionRun(ctx context.Context, id ActionRunID, token string, status ActionRunStatus, output json.RawMessage, errorText string) error {
	now := time.Now().UTC()
	update := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(xid.ID(id)), entactionrun.LeaseTokenEQ(token), entactionrun.StatusEQ(entactionrun.StatusRunning),
	).SetStatus(entactionrun.Status(status.String())).
		SetFinishedAt(now).
		ClearLeaseToken().
		ClearLeaseExpiresAt()
	if len(output) > 0 {
		update.SetOutput(output)
	}
	if errorText != "" {
		update.SetErrorText(errorText)
	}
	updated, err := update.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if updated != 1 {
		return ErrActionRunNotFound
	}
	return r.AggregateRun(ctx, id)
}

func (r *Repository) CancelQueuedActionRun(ctx context.Context, id ActionRunID) (bool, error) {
	now := time.Now().UTC()
	updated, err := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(xid.ID(id)), entactionrun.StatusEQ(entactionrun.StatusQueued),
	).SetStatus(entactionrun.StatusCancelled).SetFinishedAt(now).Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	if updated == 1 {
		return true, r.AggregateRun(ctx, id)
	}
	return false, nil
}

func (r *Repository) FinishCancellation(ctx context.Context, id ActionRunID) error {
	now := time.Now().UTC()
	updated, err := r.db.TrailActionRun.Update().Where(
		entactionrun.IDEQ(xid.ID(id)),
		entactionrun.StatusIn(entactionrun.StatusQueued, entactionrun.StatusRunning),
	).SetStatus(entactionrun.StatusCancelled).
		SetFinishedAt(now).
		ClearLeaseToken().
		ClearLeaseExpiresAt().
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if updated == 0 {
		return ErrActionRunNotFound
	}
	return r.AggregateRun(ctx, id)
}

func (r *Repository) AggregateRun(ctx context.Context, actionRunID ActionRunID) error {
	row, err := r.db.TrailActionRun.Get(ctx, xid.ID(actionRunID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	actions, err := r.db.TrailActionRun.Query().Where(entactionrun.RunIDEQ(row.RunID)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	terminal := true
	allCancelled := len(actions) > 0
	attention := false
	hasRunning := false
	hasQueued := false
	for _, action := range actions {
		switch action.Status {
		case entactionrun.StatusQueued:
			hasQueued, terminal = true, false
		case entactionrun.StatusRunning:
			hasRunning, terminal = true, false
		case entactionrun.StatusBlocked, entactionrun.StatusFailed:
			attention, allCancelled = true, false
		case entactionrun.StatusCompleted:
			allCancelled = false
		case entactionrun.StatusCancelled:
		}
	}
	if !terminal {
		status := entrun.StatusQueued
		if hasRunning {
			status = entrun.StatusRunning
		} else if !hasQueued {
			status = entrun.StatusRunning
		}
		_, err = r.db.TrailRun.UpdateOneID(row.RunID).SetStatus(status).ClearFinishedAt().Save(ctx)
		return fault.Wrap(err, fctx.With(ctx))
	}
	status := entrun.StatusCompleted
	if attention {
		status = entrun.StatusAttentionRequired
	} else if allCancelled {
		status = entrun.StatusCancelled
	}
	_, err = r.db.TrailRun.UpdateOneID(row.RunID).SetStatus(status).SetFinishedAt(time.Now().UTC()).Save(ctx)
	return fault.Wrap(err, fctx.With(ctx))
}

func (r *Repository) RecordAttentionNotification(ctx context.Context, id ActionRunID, creatorID xid.ID, target string) (bool, error) {
	recorded := false
	err := ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		updated, err := tx.TrailActionRun.Update().Where(
			entactionrun.IDEQ(xid.ID(id)), entactionrun.NotifiedAtIsNil(),
		).SetNotifiedAt(time.Now().UTC()).Save(ctx)
		if err != nil {
			return err
		}
		if updated == 0 {
			return nil
		}
		if _, err := tx.Notification.Create().
			SetOwnerID(creatorID).
			SetEventType("trail_run_attention").
			SetTarget(target).
			SetRead(false).
			Save(ctx); err != nil {
			return err
		}
		recorded = true
		return nil
	})
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return recorded, nil
}

func mapTrail(row *ent.Trail) (*Trail, error) {
	creator, err := account.MapRef(row.Edges.Creator)
	if err != nil {
		return nil, err
	}
	status, err := NewStatus(row.Status.String())
	if err != nil {
		return nil, err
	}
	trigger, err := decodeTrigger(row.TriggerType, row.TriggerConfig)
	if err != nil {
		return nil, err
	}
	result := &Trail{
		ID:               ID(row.ID),
		Creator:          *creator,
		Name:             row.Name,
		Description:      row.Description,
		Status:           status,
		Trigger:          trigger,
		NextOccurrenceAt: row.NextOccurrenceAt,
		LastOccurrenceAt: row.LastOccurrenceAt,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
	for _, action := range row.Edges.Actions {
		mapped, err := mapAction(action)
		if err != nil {
			return nil, err
		}
		result.Actions = append(result.Actions, mapped)
	}
	return result, nil
}

func mapAction(row *ent.TrailAction) (*Action, error) {
	kind, err := NewActionKind(row.Kind)
	if err != nil {
		return nil, err
	}
	return &Action{
		ID:         ActionID(row.ID),
		TrailID:    ID(row.TrailID),
		Kind:       kind,
		Position:   row.Position,
		Config:     row.Config,
		ArchivedAt: row.ArchivedAt,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func mapRun(row *ent.TrailRun) (*Run, error) {
	kind, err := NewRunKind(row.Kind.String())
	if err != nil {
		return nil, err
	}
	status, err := NewRunStatus(row.Status.String())
	if err != nil {
		return nil, err
	}
	var trigger *TriggerEvent
	if decoded, err := decodeTriggerEvent(row.TriggerPayload); err == nil {
		trigger = &decoded
	}
	result := &Run{
		ID:           RunID(row.ID),
		TrailID:      ID(row.TrailID),
		InitiatedBy:  row.InitiatedByID,
		Kind:         kind,
		Trigger:      trigger,
		ScheduledFor: row.ScheduledFor,
		Status:       status,
		FinishedAt:   row.FinishedAt,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	for _, action := range row.Edges.ActionRuns {
		mapped, err := mapActionRun(action)
		if err != nil {
			return nil, err
		}
		result.ActionRuns = append(result.ActionRuns, mapped)
	}
	return result, nil
}

func mapActionRun(row *ent.TrailActionRun) (*ActionRun, error) {
	kind, err := NewActionKind(row.Kind)
	if err != nil {
		return nil, err
	}
	status, err := NewActionRunStatus(row.Status.String())
	if err != nil {
		return nil, err
	}
	result := &ActionRun{
		ID:             ActionRunID(row.ID),
		RunID:          RunID(row.RunID),
		ActionID:       ActionID(row.ActionID),
		Kind:           kind,
		Config:         row.Config,
		Status:         status,
		LeaseToken:     row.LeaseToken,
		LeaseExpiresAt: row.LeaseExpiresAt,
		Output:         row.Output,
		Target:         row.Target,
		ErrorText:      row.ErrorText,
		StartedAt:      row.StartedAt,
		FinishedAt:     row.FinishedAt,
		NotifiedAt:     row.NotifiedAt,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
	if row.Edges.Action != nil {
		result.Action, err = mapAction(row.Edges.Action)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ptrTime(value time.Time) *time.Time { return &value }
