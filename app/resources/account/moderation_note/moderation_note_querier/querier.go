package moderation_note_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/moderation_note"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/internal/ent"
	entmoderationnote "github.com/Southclaws/storyden/internal/ent/moderationnote"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (q *Querier) ListByAccountID(ctx context.Context, accountID xid.ID) (moderation_note.Notes, error) {
	notes, err := q.db.ModerationNote.Query().
		Where(entmoderationnote.AccountIDEQ(accountID)).
		WithAuthor().
		Order(ent.Desc(entmoderationnote.FieldCreatedAt), ent.Desc(entmoderationnote.FieldID)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(notes))
	for _, note := range notes {
		roleTargets = append(roleTargets, note.Edges.Author)
	}
	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(notes, moderation_note.Map)
}
