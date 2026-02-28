package account_repo

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	entpredicate "github.com/Southclaws/storyden/internal/ent/predicate"
	role_ent "github.com/Southclaws/storyden/internal/ent/role"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
)

func (r *Repository) GetByID(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error) {
	ctx, span := r.ins.Instrument(ctx, kv.String("account_id", id.String()))
	defer span.End()

	q := r.db.Account.
		Query().
		Where(account_ent.ID(xid.ID(id))).
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator()
		}).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, roleHydrationTargets(result)...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapAccount(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.accountCache.storeAccount(ctx, &acc.Account)

	return acc, nil
}

func (r *Repository) LookupByHandle(ctx context.Context, handle string) (*account.AccountWithEdges, bool, error) {
	ctx, span := r.ins.Instrument(ctx, kv.String("handle", handle))
	defer span.End()

	q := r.db.Account.
		Query().
		Where(account_ent.Handle(handle)).
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator()
		}).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, roleHydrationTargets(result)...); err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapAccount(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	r.accountCache.storeAccount(ctx, &acc.Account)

	return acc, true, nil
}

func (r *Repository) GetRefByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	ctx, span := r.ins.Instrument(ctx, kv.String("account_id", id.String()))
	defer span.End()

	if cached, ok := r.accountCache.get(ctx, xid.ID(id)); ok {
		ctx = span.Annotate(kv.Bool("cache_hit", true))
		if err := r.refreshCachedRoles(ctx, cached); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return cached, nil
	}
	ctx = span.Annotate(kv.Bool("cache_hit", false))

	result, err := r.db.Account.
		Query().
		Where(account_ent.ID(xid.ID(id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, result); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapRef(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r.accountCache.storeAccount(ctx, acc)

	return acc, nil
}

func (r *Repository) refreshCachedRoles(ctx context.Context, acc *account.Account) error {
	ctx, span := r.ins.Instrument(ctx, kv.String("account_id", acc.ID.String()))
	defer span.End()

	hydrationTarget := &ent.Account{
		ID:        xid.ID(acc.ID),
		CreatedAt: acc.CreatedAt,
		Admin:     acc.Admin,
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, hydrationTarget); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := held.MapList(hydrationTarget.Edges.AccountRoles)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc.Roles = roles

	return nil
}

func (r *Repository) ProbeMany(ctx context.Context, handles ...string) ([]*account.Account, error) {
	ctx, span := r.ins.Instrument(ctx, kv.Int("handles_count", len(handles)))
	defer span.End()

	if len(handles) == 0 {
		return []*account.Account{}, nil
	}

	accounts, err := r.db.Account.
		Query().
		Where(
			account_ent.HandleIn(handles...),
			account_ent.DeletedAtIsNil(),
		).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, accounts...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(accounts, account.MapRef)
}

func (r *Repository) ListByHeldPermission(ctx context.Context, perms ...rbac.Permission) ([]*account.Account, error) {
	ctx, span := r.ins.Instrument(ctx, kv.Int("permissions_count", len(perms)))
	defer span.End()

	if len(perms) == 0 {
		return []*account.Account{}, nil
	}

	predicates := make([]entpredicate.Account, 0, len(perms)+1)
	predicates = append(predicates, account_ent.Admin(true))

	for _, perm := range perms {
		p := perm.String()
		predicates = append(predicates, account_ent.HasRolesWith(entpredicate.Role(func(s *sql.Selector) {
			s.Where(sqljson.ValueContains(role_ent.FieldPermissions, p))
		})))
	}

	accounts, err := r.db.Account.Query().
		Where(account_ent.DeletedAtIsNil()).
		Where(account_ent.Or(predicates...)).
		Unique(true).
		Order(account_ent.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.roleHydrator.HydrateRoleEdges(ctx, accounts...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(accounts, account.MapRef)
}

func roleHydrationTargets(acc *ent.Account) []*ent.Account {
	targets := []*ent.Account{acc}

	if invitedBy := acc.Edges.InvitedBy; invitedBy != nil {
		creator, err := invitedBy.Edges.CreatorOrErr()
		if err == nil {
			targets = append(targets, creator)
		}
	}

	return targets
}
