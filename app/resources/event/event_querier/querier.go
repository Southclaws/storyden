package event_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) Probe(ctx context.Context, mark event_ref.QueryKey) (*event_ref.Event, error) {
	query := q.db.Event.Query()

	query.Where(mark.Predicate())

	r, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return event_ref.Map(r)
}

func (q *Querier) Get(ctx context.Context, mark event_ref.QueryKey) (*event.Event, error) {
	query := q.db.Event.Query()

	query.Where(mark.Predicate())

	query.WithParticipants(func(epq *ent.EventParticipantQuery) {
		epq.WithAccount()
	})
	query.WithPrimaryImage()
	query.WithThread(func(pq *ent.PostQuery) {
		pq.WithCategory()
		pq.WithAuthor()
		pq.WithPosts()
		pq.WithReacts()
	})

	r, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err := event.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return evt, nil
}

func (q *Querier) List(ctx context.Context) ([]*event_ref.Event, error) {
	query := q.db.Event.Query()

	query.WithParticipants()
	query.WithPrimaryImage()

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evts, err := dt.MapErr(r, event_ref.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return evts, nil
}
