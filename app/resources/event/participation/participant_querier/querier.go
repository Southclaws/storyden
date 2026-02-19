package participant_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/internal/ent"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (w *Querier) Lookup(ctx context.Context, mk event_ref.QueryKey, accountID account.AccountID) (*participation.EventParticipant, bool, error) {
	ep, err := w.db.EventParticipant.Query().
		Where(mk.ParticipantPredicate()).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := w.roleQuerier.BuildSingleHydrator(ctx, ep.Edges.Account)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	participant, err := participation.Map(ep, roleHydrator.Hydrate)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return participant, true, nil
}
