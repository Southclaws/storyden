package event_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (q *Querier) Probe(ctx context.Context, mark event_ref.QueryKey) (*event_ref.Event, error) {
	query := q.db.Event.Query()

	query.Where(mark.Predicate())

	r, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return event_ref.Map(r, roleHydrator.Hydrate)
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

	accounts := make([]*ent.Account, 0, len(r.Edges.Participants)+1)
	if r.Edges.Thread != nil && r.Edges.Thread.Edges.Author != nil {
		accounts = append(accounts, r.Edges.Thread.Edges.Author)
	}
	for _, p := range r.Edges.Participants {
		if p != nil && p.Edges.Account != nil {
			accounts = append(accounts, p.Edges.Account)
		}
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, accounts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err := event.Map(r, roleHydrator.Hydrate)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return evt, nil
}

func (q *Querier) List(ctx context.Context) ([]*event_ref.Event, error) {
	query := q.db.Event.Query()

	query.WithParticipants(func(epq *ent.EventParticipantQuery) {
		epq.WithAccount()
	})
	query.WithPrimaryImage()

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accounts := make([]*ent.Account, 0)
	for _, evt := range r {
		for _, participant := range evt.Edges.Participants {
			if participant != nil && participant.Edges.Account != nil {
				accounts = append(accounts, participant.Edges.Account)
			}
		}
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, accounts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evts, err := dt.MapErr(r, func(in *ent.Event) (*event_ref.Event, error) {
		return event_ref.Map(in, roleHydrator.Hydrate)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return evts, nil
}
