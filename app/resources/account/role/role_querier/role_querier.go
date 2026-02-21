package role_querier

import (
	"context"
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

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) Get(ctx context.Context, id role.RoleID) (*role.Role, error) {
	r, err := q.db.Role.Get(ctx, xid.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rl, nil
}

func (q *Querier) List(ctx context.Context) (role.Roles, error) {
	roles, err := q.db.Role.Query().Where(ent_role.IDNotIn(
		xid.ID(role.DefaultRoleGuestID),
		xid.ID(role.DefaultRoleMemberID),
		xid.ID(role.DefaultRoleAdminID),
	)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.MapList(roles)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defaultRole, err := q.GetMemberRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	guestRole, err := q.GetGuestRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	adminRole, err := q.GetAdminRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped = append(mapped, defaultRole, guestRole, adminRole)

	sort.Sort(mapped)

	return mapped, nil
}

func (q *Querier) ListFor(ctx context.Context, account *ent.Account) (held.Roles, error) {
	roles, err := q.db.AccountRoles.
		Query().
		Where(
			ent_account_role.AccountID(account.ID),
			ent_account_role.HasRoleWith(ent_role.IDNotIn(
				xid.ID(role.DefaultRoleGuestID),
				xid.ID(role.DefaultRoleMemberID),
				xid.ID(role.DefaultRoleAdminID),
			)),
		).
		WithRole(func(rq *ent.RoleQuery) {
			rq.Order(ent.Asc(ent_role.FieldSortKey))
		}).
		Order(ent.Asc(ent_account_role.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := held.MapList(roles, account.Admin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, memberRole, _, err := q.lookupDefaultRoles(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// If the default member role has not been modified (aka not added to the DB
	// with custom permissions) we add the default manually.
	if memberRole != nil {
		defaultRole, err := role.Map(memberRole)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		mapped = append(mapped, &held.Role{
			Role:     *defaultRole,
			Assigned: account.CreatedAt,
			Badge:    false,
			Default:  true,
		})
	} else {
		mapped = append(mapped, &held.Role{
			Role: role.DefaultRoleMember,
		})
	}

	sort.Sort(mapped)

	return mapped, nil
}

func (q *Querier) GetMemberRole(ctx context.Context) (*role.Role, error) {
	_, memberRole, _, err := q.lookupDefaultRoles(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if memberRole == nil {
		return &role.DefaultRoleMember, nil
	}

	return role.Map(memberRole)
}

func (q *Querier) GetGuestRole(ctx context.Context) (*role.Role, error) {
	guestRole, _, _, err := q.lookupDefaultRoles(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if guestRole == nil {
		return &role.DefaultRoleGuest, nil
	}

	return role.Map(guestRole)
}

func (q *Querier) GetAdminRole(ctx context.Context) (*role.Role, error) {
	_, _, adminRole, err := q.lookupDefaultRoles(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if adminRole == nil {
		return &role.DefaultRoleAdmin, nil
	}

	return role.Map(adminRole)
}

func (q *Querier) lookupDefaultRoles(ctx context.Context) (*ent.Role, *ent.Role, *ent.Role, error) {
	roles, err := q.db.Role.Query().Where(ent_role.IDIn(
		xid.ID(role.DefaultRoleGuestID),
		xid.ID(role.DefaultRoleMemberID),
		xid.ID(role.DefaultRoleAdminID),
	)).All(ctx)
	if err != nil {
		return nil, nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	var guestRole *ent.Role
	var memberRole *ent.Role
	var adminRole *ent.Role

	for _, r := range roles {
		if r.ID == xid.ID(role.DefaultRoleGuestID) {
			guestRole = r
		} else if r.ID == xid.ID(role.DefaultRoleMemberID) {
			memberRole = r
		} else if r.ID == xid.ID(role.DefaultRoleAdminID) {
			adminRole = r
		}
	}

	return guestRole, memberRole, adminRole, nil
}
