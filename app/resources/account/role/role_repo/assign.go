package role_repo

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
)

func splitAssignMutations(mutations ...role_assign.Mutation) (adds, removes []xid.ID, admin opt.Optional[bool]) {
	for _, m := range mutations {
		switch {
		case m.IsDelete() && m.ID() == role.DefaultRoleAdminID:
			admin = opt.New(false)
		case m.IsDelete():
			removes = append(removes, xid.ID(m.ID()))
		case m.ID() == role.DefaultRoleAdminID:
			admin = opt.New(true)
		default:
			adds = append(adds, xid.ID(m.ID()))
		}
	}

	return
}

func (h *Repository) UpdateRoles(ctx context.Context, accountID xid.ID, roles ...role_assign.Mutation) error {
	update := h.db.Account.UpdateOneID(accountID)
	mutation := update.Mutation()

	filtered := make([]role_assign.Mutation, 0, len(roles))
	for _, m := range roles {
		if m.ID() == role.DefaultRoleMemberID {
			continue
		}

		filtered = append(filtered, m)
	}

	adds, removes, admin := splitAssignMutations(filtered...)

	mutation.AddRoleIDs(adds...)
	mutation.RemoveRoleIDs(removes...)
	if a, ok := admin.Get(); ok {
		mutation.SetAdmin(a)
	}

	_, err := update.Save(ctx)
	if err != nil {
		if !ent.IsConstraintError(err) {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	relationships, err := h.db.AccountRoles.Query().Where(
		ent_account_role.AccountID(accountID),
		ent_account_role.RoleIDNotIn(
			xid.ID(role.DefaultRoleGuestID),
			xid.ID(role.DefaultRoleMemberID),
			xid.ID(role.DefaultRoleAdminID),
		),
	).Order(ent.Asc(ent_account_role.FieldCreatedAt)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	assignments := mapAssignmentsFromRelationships(relationships)

	if err := h.storeAssignments(ctx, accountID, assignments); err != nil {
		if recoveryErr := h.recoverFromCacheWriteFailure(ctx, "assign.update_roles.store_assignments", err); recoveryErr != nil {
			return fault.Wrap(recoveryErr, fctx.With(ctx))
		}
	}

	return nil
}
