package role_repo

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

func (h *Repository) UpdateBadge(ctx context.Context, accountID xid.ID, roleID role.RoleID, badge bool) error {
	tx, err := h.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		_ = tx.Rollback()
	}()

	filters := []predicate.AccountRoles{
		accountroles.AccountIDEQ(accountID),
		accountroles.RoleIDEQ(xid.ID(roleID)),
	}

	// Only one role can be set as a badge, clear all first, then set if true.
	if err := tx.AccountRoles.Update().
		Where(accountroles.AccountIDEQ(accountID)).
		ClearBadge().
		Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if badge {
		affected, err := tx.AccountRoles.Update().
			Where(filters...).
			SetBadge(true).
			Save(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		if affected == 0 {
			return fault.Wrap(fault.New("role assignment not found"), fctx.With(ctx), ftag.With(ftag.NotFound))
		}
	}

	relationships, err := tx.AccountRoles.Query().Where(
		accountroles.AccountIDEQ(accountID),
		accountroles.RoleIDNotIn(
			xid.ID(role.DefaultRoleGuestID),
			xid.ID(role.DefaultRoleMemberID),
			xid.ID(role.DefaultRoleAdminID),
		),
	).Order(ent.Asc(accountroles.FieldCreatedAt)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	assignments := mapAssignmentsFromRelationships(relationships)

	if err := h.storeAssignments(ctx, accountID, assignments); err != nil {
		if recoveryErr := h.recoverFromCacheWriteFailure(ctx, "badge.update_badge.store_assignments", err); recoveryErr != nil {
			slog.Error("role repository cache write recovery failed after committed badge update",
				slog.String("operation", "badge.update_badge.store_assignments"),
				slog.String("cache_write_error", err.Error()),
				slog.String("recovery_error", recoveryErr.Error()))
		}
	}

	return nil
}
