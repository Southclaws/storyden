package role_repo

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/ent"
)

func (h *Repository) BuildSingleHydrator(ctx context.Context, account *ent.Account) (Hydrator, error) {
	roles, err := h.ListFor(ctx, account)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return BuildSingleHydrator(account.ID, roles), nil
}

func (h *Repository) BuildMultiHydrator(ctx context.Context, accounts []*ent.Account) (Hydrator, error) {
	roleLookup, err := h.ListForMany(ctx, accounts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return BuildMultiHydrator(roleLookup), nil
}
