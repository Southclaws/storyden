package role_hydrate

import (
	"context"
	"math"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Hydrator struct {
	ins                spanner.Instrumentation
	roleRepo           *role_repo.Repository
	assignmentResolver *role_assign.Assignment
}

func New(
	ins spanner.Builder,
	roleRepo *role_repo.Repository,
	assignmentResolver *role_assign.Assignment,
) *Hydrator {
	return &Hydrator{
		ins:                ins.Build(),
		roleRepo:           roleRepo,
		assignmentResolver: assignmentResolver,
	}
}

func (r *Hydrator) HydrateRoleEdges(ctx context.Context, accounts ...*ent.Account) error {
	ctx, span := r.ins.Instrument(ctx, kv.Int("accounts_count", len(accounts)))
	defer span.End()

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

	ctx = span.Annotate(kv.Int("deduped_accounts_count", len(accountIDs)))

	idsByAccount, byAccountRows, err := r.assignmentResolver.ResolveRoleIDs(ctx, accountIDs)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	roleByID, err := r.rolesByID(ctx, idsByAccount, hasAdmin)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	staleAccountIDs := findAccountsWithMissingRoles(idsByAccount, roleByID)
	if len(staleAccountIDs) > 0 {
		ctx = span.Annotate(kv.Int("stale_accounts_count", len(staleAccountIDs)))
		refreshedIDs, refreshedRows, err := r.assignmentResolver.ResolveRoleIDsFresh(ctx, staleAccountIDs)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		for _, accountID := range staleAccountIDs {
			idsByAccount[accountID] = refreshedIDs[accountID]
			if rows, ok := refreshedRows[accountID]; ok {
				byAccountRows[accountID] = rows
			} else {
				delete(byAccountRows, accountID)
			}
		}

		roleByID, err = r.rolesByID(ctx, idsByAccount, hasAdmin)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	roleLess := func(a, b *ent.AccountRoles) bool {
		roleKey := func(in *ent.AccountRoles) (float64, string) {
			if in == nil {
				return math.MaxFloat64, ""
			}

			rl := roleByID[role.RoleID(in.RoleID)]
			if rl == nil {
				return math.MaxFloat64, in.RoleID.String()
			}

			return roleSortKey(rl.ID, rl.SortKey), rl.ID.String()
		}

		sa, ia := roleKey(a)
		sb, ib := roleKey(b)
		if sa == sb {
			return ia < ib
		}

		return sa > sb
	}

	memberRole := roleByID[role.DefaultRoleMemberID]
	if memberRole == nil {
		copy := role.DefaultRoleMember
		memberRole = &copy
	}
	memberRoleEnt := mapRoleToEnt(memberRole)

	adminRole := roleByID[role.DefaultRoleAdminID]
	if adminRole == nil {
		copy := role.DefaultRoleAdmin
		adminRole = &copy
	}
	adminRoleEnt := mapRoleToEnt(adminRole)

	for _, acc := range accounts {
		if acc == nil {
			continue
		}

		heldRoles := make([]*ent.AccountRoles, 0, len(idsByAccount[acc.ID])+2)

		if rows, ok := byAccountRows[acc.ID]; ok {
			for _, ar := range rows {
				roleEdge := roleByID[role.RoleID(ar.RoleID)]
				if roleEdge == nil {
					continue
				}

				ar.Edges.Role = mapRoleToEnt(roleEdge)
				heldRoles = append(heldRoles, ar)
			}
		} else {
			for _, roleID := range idsByAccount[acc.ID] {
				roleEdge := roleByID[role.RoleID(roleID)]
				if roleEdge == nil {
					continue
				}

				heldRoles = append(heldRoles, &ent.AccountRoles{
					ID:        xid.New(),
					CreatedAt: acc.CreatedAt,
					AccountID: acc.ID,
					RoleID:    roleID,
					Edges: ent.AccountRolesEdges{
						Role: mapRoleToEnt(roleEdge),
					},
				})
			}
		}

		acc.Edges.AccountRoles = heldRoles
		appendDefaultRoleEdge(acc, memberRoleEnt)
		if acc.Admin {
			appendDefaultRoleEdge(acc, adminRoleEnt)
		}

		sort.SliceStable(acc.Edges.AccountRoles, func(i, j int) bool {
			return roleLess(acc.Edges.AccountRoles[i], acc.Edges.AccountRoles[j])
		})
	}

	return nil
}

func mapRoleToEnt(in *role.Role) *ent.Role {
	if in == nil {
		return nil
	}

	return &ent.Role{
		ID:          xid.ID(in.ID),
		Name:        in.Name,
		Colour:      in.Colour,
		Permissions: dt.Map(in.Permissions.List(), func(p rbac.Permission) string { return p.String() }),
		SortKey:     in.SortKey,
		Metadata:    in.Metadata,
	}
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

func (r *Hydrator) rolesByID(ctx context.Context, idsByAccount map[xid.ID][]xid.ID, hasAdmin bool) (map[role.RoleID]*role.Role, error) {
	ctx, span := r.ins.Instrument(ctx,
		kv.Int("accounts_count", len(idsByAccount)),
		kv.Bool("has_admin", hasAdmin),
	)
	defer span.End()

	roleIDSet := map[role.RoleID]struct{}{
		role.DefaultRoleMemberID: {},
	}
	if hasAdmin {
		roleIDSet[role.DefaultRoleAdminID] = struct{}{}
	}
	for _, roleIDs := range idsByAccount {
		for _, roleID := range roleIDs {
			roleIDSet[role.RoleID(roleID)] = struct{}{}
		}
	}

	roleIDs := make([]role.RoleID, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}

	ctx = span.Annotate(kv.Int("role_ids_count", len(roleIDs)))

	roleByID, err := r.roleRepo.GetMany(ctx, roleIDs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return roleByID, nil
}

func findAccountsWithMissingRoles(idsByAccount map[xid.ID][]xid.ID, roleByID map[role.RoleID]*role.Role) []xid.ID {
	stale := make([]xid.ID, 0)
	for accountID, roleIDs := range idsByAccount {
		for _, roleID := range roleIDs {
			if roleByID[role.RoleID(roleID)] == nil {
				stale = append(stale, accountID)
				break
			}
		}
	}
	return stale
}
