package role_repo

import (
	"context"
	"log/slog"
	"sort"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

func (h *Repository) Get(ctx context.Context, id role.RoleID) (*role.Role, error) {
	rl, _, err := h.getRole(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rl, nil
}

func (h *Repository) GetGuestRole(ctx context.Context) (*role.Role, error) {
	return h.Get(ctx, role.DefaultRoleGuestID)
}

func (h *Repository) GetMemberRole(ctx context.Context) (*role.Role, error) {
	return h.Get(ctx, role.DefaultRoleMemberID)
}

func (h *Repository) GetAdminRole(ctx context.Context) (*role.Role, error) {
	return h.Get(ctx, role.DefaultRoleAdminID)
}

func (h *Repository) List(ctx context.Context) (role.Roles, error) {
	ids, ok := h.cachedCustomRoleOrdering(ctx)
	if !ok {
		roles, err := h.listFromDB(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		slog.Warn("role repository read-through fallback used for custom role ordering")
		return roles, nil
	}

	custom := make(role.Roles, 0, len(ids))
	for _, id := range ids {
		rl, _, err := h.getRole(ctx, id)
		if err != nil {
			if ent.IsNotFound(err) {
				// Role list is out of sync, rebuild immediately from DB.
				roles, ferr := h.listFromDB(ctx)
				if ferr != nil {
					return nil, fault.Wrap(ferr, fctx.With(ctx))
				}

				slog.Warn("role repository read-through fallback used for stale custom role ordering")
				return roles, nil
			}

			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		custom = append(custom, rl)
	}

	withDefaults, err := h.appendDefaultRoles(ctx, custom)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return withDefaults, nil
}

func (h *Repository) ListFor(ctx context.Context, acc *ent.Account) (held.Roles, error) {
	assignments, ok := h.cachedAssignments(ctx, xid.ID(acc.ID))
	cacheMiss := !ok
	staleRebuild := false

	if ok {
		roles, stale, err := h.hydrateHeld(ctx, acc, assignments)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if !stale {
			return roles, nil
		}

		staleRebuild = true
	}

	assignments, err := h.assignmentsFromDB(ctx, xid.ID(acc.ID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_ = h.storeAssignments(ctx, xid.ID(acc.ID), assignments)

	roles, _, err := h.hydrateHeld(ctx, acc, assignments)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if staleRebuild {
		slog.Warn("role repository read-through fallback used for stale account role assignments",
			slog.String("account_id", acc.ID.String()))
	} else if cacheMiss {
		slog.Warn("role repository read-through fallback used for account role assignments",
			slog.String("account_id", acc.ID.String()))
	}

	return roles, nil
}

func (h *Repository) ListForMany(ctx context.Context, accounts []*ent.Account) (map[xid.ID]held.Roles, error) {
	uniqueAccounts := map[xid.ID]*ent.Account{}
	for _, acc := range accounts {
		if acc == nil {
			continue
		}

		uniqueAccounts[acc.ID] = acc
	}

	if len(uniqueAccounts) == 0 {
		return map[xid.ID]held.Roles{}, nil
	}

	assignmentsByAccount := map[xid.ID][]cachedAssignment{}
	missing := make([]xid.ID, 0, len(uniqueAccounts))

	for id := range uniqueAccounts {
		assignments, ok := h.cachedAssignments(ctx, id)
		if ok {
			assignmentsByAccount[id] = assignments
			continue
		}

		missing = append(missing, id)
	}

	if len(missing) > 0 {
		relationships, err := h.db.AccountRoles.Query().Where(
			ent_account_role.AccountIDIn(missing...),
			ent_account_role.RoleIDNotIn(
				xid.ID(role.DefaultRoleGuestID),
				xid.ID(role.DefaultRoleMemberID),
				xid.ID(role.DefaultRoleAdminID),
			),
		).Order(ent.Asc(ent_account_role.FieldCreatedAt)).All(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		for _, relationship := range relationships {
			assignmentsByAccount[relationship.AccountID] = append(assignmentsByAccount[relationship.AccountID], cachedAssignment{
				RoleID:     relationship.RoleID.String(),
				AssignedAt: relationship.CreatedAt,
				Badge:      relationship.Badge != nil && *relationship.Badge,
			})
		}

		for _, id := range missing {
			_ = h.storeAssignments(ctx, id, assignmentsByAccount[id])
		}

		slog.Warn("role repository read-through fallback used for account role assignments (batch)",
			slog.Int("missing_count", len(missing)))
	}

	out := make(map[xid.ID]held.Roles, len(uniqueAccounts))

	for id, acc := range uniqueAccounts {
		assignments := assignmentsByAccount[id]
		roles, stale, err := h.hydrateHeld(ctx, acc, assignments)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if stale {
			refetched, err := h.assignmentsFromDB(ctx, id)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			_ = h.storeAssignments(ctx, id, refetched)

			roles, _, err = h.hydrateHeld(ctx, acc, refetched)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			slog.Warn("role repository read-through fallback used for stale account role assignments",
				slog.String("account_id", id.String()))
		}

		out[id] = roles
	}

	return out, nil
}

func (h *Repository) listFromDB(ctx context.Context) (role.Roles, error) {
	customRoles, err := h.db.Role.Query().Where(ent_role.IDNotIn(
		xid.ID(role.DefaultRoleGuestID),
		xid.ID(role.DefaultRoleMemberID),
		xid.ID(role.DefaultRoleAdminID),
	)).Order(ent.Asc(ent_role.FieldSortKey)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.MapList(customRoles)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ids := make([]role.RoleID, 0, len(mapped))
	for _, rl := range mapped {
		ids = append(ids, rl.ID)
		_ = h.storeRole(ctx, rl, true)
	}

	_ = h.storeCustomRoleOrdering(ctx, ids)

	withDefaults, err := h.appendDefaultRoles(ctx, mapped)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return withDefaults, nil
}

func (h *Repository) appendDefaultRoles(ctx context.Context, custom role.Roles) (role.Roles, error) {
	guestRole, _, err := h.getRole(ctx, role.DefaultRoleGuestID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	memberRole, _, err := h.getRole(ctx, role.DefaultRoleMemberID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	adminRole, _, err := h.getRole(ctx, role.DefaultRoleAdminID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all := make(role.Roles, 0, len(custom)+3)
	all = append(all, custom...)
	all = append(all, memberRole, guestRole, adminRole)
	sort.Sort(all)

	return all, nil
}

func (h *Repository) getRole(ctx context.Context, id role.RoleID) (*role.Role, bool, error) {
	if cached, ok := h.cachedRole(ctx, id); ok {
		parsed, err := cached.toRole()
		if err == nil {
			return parsed, cached.Persisted, nil
		}
	}

	r, err := h.db.Role.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			defaultRole, ok := defaultRole(id)
			if !ok {
				return nil, false, fault.Wrap(err, fctx.With(ctx))
			}

			defaultCopy := defaultRole
			_ = h.storeRole(ctx, &defaultCopy, false)
			return &defaultCopy, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.Map(r)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	_ = h.storeRole(ctx, mapped, true)
	slog.Warn("role repository read-through fallback used for role data",
		slog.String("role_id", id.String()))

	return mapped, true, nil
}
