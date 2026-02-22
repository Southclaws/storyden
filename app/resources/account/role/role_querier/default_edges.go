package role_querier

import (
	"context"
	"math"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

func (q *Querier) HydrateRoleEdges(ctx context.Context, accounts ...*ent.Account) error {
	accountIDs := make([]xid.ID, 0, len(accounts))
	seen := make(map[xid.ID]struct{}, len(accounts))
	hasAdmin := false
	for _, acc := range accounts {
		if acc == nil {
			continue
		}

		if acc.Admin {
			hasAdmin = true
		}

		if _, ok := seen[acc.ID]; ok {
			continue
		}
		seen[acc.ID] = struct{}{}
		accountIDs = append(accountIDs, acc.ID)
	}
	if len(accountIDs) == 0 {
		return nil
	}

	accountRoles, err := q.db.AccountRoles.Query().
		Where(ent_account_role.AccountIDIn(accountIDs...)).
		Order(ent.Asc(ent_account_role.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	roleIDs := map[xid.ID]struct{}{
		xid.ID(role.DefaultRoleMemberID): {},
	}
	if hasAdmin {
		roleIDs[xid.ID(role.DefaultRoleAdminID)] = struct{}{}
	}
	for _, held := range accountRoles {
		roleIDs[held.RoleID] = struct{}{}
	}

	roleIDList := make([]xid.ID, 0, len(roleIDs))
	for roleID := range roleIDs {
		roleIDList = append(roleIDList, roleID)
	}

	roleRows, err := q.db.Role.Query().Where(ent_role.IDIn(roleIDList...)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	roleByID := lo.KeyBy(roleRows, func(r *ent.Role) xid.ID { return r.ID })
	roleLess := func(a, b *ent.AccountRoles) bool {
		roleKey := func(in *ent.AccountRoles) (float64, string) {
			if in == nil {
				return math.MaxFloat64, ""
			}

			rl := roleByID[in.RoleID]
			if rl == nil {
				return math.MaxFloat64, in.RoleID.String()
			}

			return roleSortKey(role.RoleID(rl.ID), rl.SortKey), rl.ID.String()
		}

		sa, ia := roleKey(a)
		sb, ib := roleKey(b)
		if sa == sb {
			return ia < ib
		}

		return sa < sb
	}

	byAccount := lo.GroupBy(accountRoles, func(r *ent.AccountRoles) xid.ID { return r.AccountID })

	memberRole := roleByID[xid.ID(role.DefaultRoleMemberID)]
	if memberRole == nil {
		memberRole = syntheticDefaultRole(role.DefaultRoleMember)
		roleByID[memberRole.ID] = memberRole
	}

	adminRole := roleByID[xid.ID(role.DefaultRoleAdminID)]
	if adminRole == nil {
		adminRole = syntheticDefaultRole(role.DefaultRoleAdmin)
		roleByID[adminRole.ID] = adminRole
	}

	for _, acc := range accounts {
		if acc == nil {
			continue
		}

		heldRoles := make([]*ent.AccountRoles, 0, len(byAccount[acc.ID])+2)
		for _, held := range byAccount[acc.ID] {
			roleEdge := roleByID[held.RoleID]
			if roleEdge == nil {
				continue
			}
			held.Edges.Role = roleEdge
			heldRoles = append(heldRoles, held)
		}

		acc.Edges.AccountRoles = heldRoles
		appendDefaultRoleEdge(acc, memberRole)
		if acc.Admin {
			appendDefaultRoleEdge(acc, adminRole)
		}

		sort.SliceStable(acc.Edges.AccountRoles, func(i, j int) bool {
			return roleLess(acc.Edges.AccountRoles[i], acc.Edges.AccountRoles[j])
		})
	}

	return nil
}

// HydrateDefaultRoleEdges is kept as a transitional alias while call-sites are
// migrated to full role-edge hydration.
func (q *Querier) HydrateDefaultRoleEdges(ctx context.Context, accounts ...*ent.Account) error {
	return q.HydrateRoleEdges(ctx, accounts...)
}

func appendDefaultRoleEdge(acc *ent.Account, defaultRole *ent.Role) {
	if defaultRole == nil {
		return
	}

	if acc.Edges.AccountRoles == nil {
		acc.Edges.AccountRoles = []*ent.AccountRoles{}
	}

	for _, held := range acc.Edges.AccountRoles {
		if held.RoleID == defaultRole.ID {
			return
		}
	}

	acc.Edges.AccountRoles = append(acc.Edges.AccountRoles, &ent.AccountRoles{
		ID:        xid.New(),
		CreatedAt: acc.CreatedAt,
		AccountID: acc.ID,
		RoleID:    defaultRole.ID,
		Edges: ent.AccountRolesEdges{
			Role: defaultRole,
		},
	})
}

func syntheticDefaultRole(in role.Role) *ent.Role {
	perms := dt.Map(in.Permissions.List(), func(p rbac.Permission) string { return p.String() })

	return &ent.Role{
		ID:          xid.ID(in.ID),
		Name:        in.Name,
		Colour:      in.Colour,
		Permissions: perms,
		SortKey:     in.SortKey,
		Metadata:    in.Metadata,
	}
}

func roleSortKey(id role.RoleID, current float64) float64 {
	switch id {
	case role.DefaultRoleGuestID:
		return role.DefaultRoleGuest.SortKey
	case role.DefaultRoleMemberID:
		return role.DefaultRoleMember.SortKey
	case role.DefaultRoleAdminID:
		return role.DefaultRoleAdmin.SortKey
	default:
		return current
	}
}
