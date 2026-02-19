package role_repo

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
)

type cachedAssignment struct {
	RoleID     string    `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
	Badge      bool      `json:"badge"`
}

func mapAssignmentsFromRelationships(relationships []*ent.AccountRoles) []cachedAssignment {
	assignments := make([]cachedAssignment, 0, len(relationships))
	for _, relationship := range relationships {
		assignments = append(assignments, cachedAssignment{
			RoleID:     relationship.RoleID.String(),
			AssignedAt: relationship.CreatedAt,
			Badge:      relationship.Badge != nil && *relationship.Badge,
		})
	}

	return assignments
}

func (h *Repository) accountRoleKey(accountID xid.ID) string {
	return accountRoleCachePrefix + accountID.String()
}

func (h *Repository) deleteAssignments(ctx context.Context, accountID xid.ID) error {
	if err := h.store.Delete(ctx, h.accountRoleKey(accountID)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (h *Repository) storeAssignments(ctx context.Context, accountID xid.ID, assignments []cachedAssignment) error {
	raw, err := json.Marshal(assignments)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := h.store.Set(ctx, h.accountRoleKey(accountID), string(raw), cacheTTL); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (h *Repository) cachedAssignments(ctx context.Context, accountID xid.ID) ([]cachedAssignment, bool) {
	raw, err := h.store.Get(ctx, h.accountRoleKey(accountID))
	if err != nil {
		return nil, false
	}

	var out []cachedAssignment
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		_ = h.deleteAssignments(ctx, accountID)
		return nil, false
	}

	return out, true
}

func (h *Repository) hydrateHeld(ctx context.Context, acc *ent.Account, assignments []cachedAssignment) (held.Roles, bool, error) {
	mapped := make(held.Roles, 0, len(assignments)+2)

	for _, assignment := range assignments {
		roleID, err := xid.FromString(assignment.RoleID)
		if err != nil {
			return nil, true, nil
		}

		rl, _, err := h.getRole(ctx, role.RoleID(roleID))
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, true, nil
			}

			return nil, false, fault.Wrap(err, fctx.With(ctx))
		}

		mapped = append(mapped, &held.Role{
			Role:     *rl,
			Assigned: assignment.AssignedAt,
			Badge:    assignment.Badge,
		})
	}

	if acc.Admin {
		mapped = append(mapped, &held.Role{
			Role:    role.DefaultRoleAdmin,
			Default: true,
		})
	}

	memberRole, persisted, err := h.getRole(ctx, role.DefaultRoleMemberID)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	if persisted {
		mapped = append(mapped, &held.Role{
			Role:     *memberRole,
			Assigned: acc.CreatedAt,
			Badge:    false,
			Default:  true,
		})
	} else {
		mapped = append(mapped, &held.Role{
			Role: role.DefaultRoleMember,
		})
	}

	sort.Sort(mapped)

	return mapped, false, nil
}

func (h *Repository) assignmentsFromDB(ctx context.Context, accountID xid.ID) ([]cachedAssignment, error) {
	relationships, err := h.db.AccountRoles.Query().Where(
		ent_account_role.AccountID(xid.ID(accountID)),
		ent_account_role.RoleIDNotIn(
			xid.ID(role.DefaultRoleGuestID),
			xid.ID(role.DefaultRoleMemberID),
			xid.ID(role.DefaultRoleAdminID),
		),
	).Order(ent.Asc(ent_account_role.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapAssignmentsFromRelationships(relationships), nil
}
