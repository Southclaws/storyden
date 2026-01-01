package account_querier

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
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	entpredicate "github.com/Southclaws/storyden/internal/ent/predicate"
	role_ent "github.com/Southclaws/storyden/internal/ent/role"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, roleQuerier *role_querier.Querier) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (d *Querier) GetByID(ctx context.Context, id account.AccountID) (*account.AccountWithEdges, error) {
	q := d.db.Account.
		Query().
		Where(account_ent.ID(xid.ID(id))).
		WithTags().
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	hr, err := d.roleQuerier.ListFor(ctx, result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapAccount(hr)(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (d *Querier) LookupByHandle(ctx context.Context, handle string) (*account.AccountWithEdges, bool, error) {
	q := d.db.Account.
		Query().
		Where(account_ent.Handle(handle)).
		WithTags().
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	hr, err := d.roleQuerier.ListFor(ctx, result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapAccount(hr)(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

func (d *Querier) ProbeMany(ctx context.Context, handles ...string) ([]*account.Account, error) {
	if len(handles) == 0 {
		return []*account.Account{}, nil
	}

	accounts, err := d.db.Account.
		Query().
		Where(account_ent.HandleIn(handles...)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(accounts, account.MapRef)
}

func (d *Querier) ListByHeldPermission(ctx context.Context, perms ...rbac.Permission) ([]*account.Account, error) {
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

	accounts, err := d.db.Account.Query().
		Where(account_ent.DeletedAtIsNil()).
		Where(account_ent.Or(predicates...)).
		Unique(true).
		Order(account_ent.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(accounts, account.MapRef)
}
