package participant_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/event/event_querier"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/eventparticipant"
)

type Writer struct {
	db      *ent.Client
	querier *event_querier.Querier
}

func New(db *ent.Client, querier *event_querier.Querier) *Writer {
	return &Writer{db: db, querier: querier}
}

type Option func(*ent.EventParticipantMutation)

func WithRole(role participation.Role) Option {
	return func(m *ent.EventParticipantMutation) {
		m.SetRole(role.String())
	}
}

func WithStatus(status participation.Status) Option {
	return func(m *ent.EventParticipantMutation) {
		m.SetStatus(status.String())
	}
}

func (w *Writer) Add(ctx context.Context, mk event_ref.QueryKey, accountID account.AccountID, opts ...Option) error {
	evt, err := w.querier.Probe(ctx, mk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	create := w.db.EventParticipant.Create()
	mutation := create.Mutation()

	mutation.SetEventID(xid.ID(evt.ID))
	mutation.SetAccountID(xid.ID(accountID))

	for _, opt := range opts {
		opt(mutation)
	}

	err = create.OnConflictColumns(eventparticipant.FieldAccountID, eventparticipant.FieldEventID).UpdateNewValues().Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) Update(ctx context.Context, mk event_ref.QueryKey, accountID account.AccountID, opts ...Option) error {
	update := w.db.EventParticipant.Update().Where(
		mk.ParticipantPredicate(),
		eventparticipant.AccountID(xid.ID(accountID)),
	)
	mutation := update.Mutation()

	for _, opt := range opts {
		opt(mutation)
	}

	_, err := update.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) Remove(ctx context.Context, mk event_ref.QueryKey, accountID account.AccountID) error {
	delete := w.db.EventParticipant.Delete()

	delete.Where(
		mk.ParticipantPredicate(),
		eventparticipant.AccountID(xid.ID(accountID)),
	)

	_, err := delete.Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
