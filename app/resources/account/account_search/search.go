package account_search

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/sortrule"
	"github.com/Southclaws/storyden/app/resources/timerange"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	auth_ent "github.com/Southclaws/storyden/internal/ent/authentication"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
	"github.com/Southclaws/storyden/internal/ent/invitation"
	entpredicate "github.com/Southclaws/storyden/internal/ent/predicate"
)

type Filter func(*ent.AccountQuery)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func WithQuery(q string) Filter {
	return func(aq *ent.AccountQuery) {
		predicates := []entpredicate.Account{
			account_ent.NameContainsFold(q),
			account_ent.HandleContainsFold(q),
			account_ent.HasEmailsWith(email_ent.EmailAddressContainsFold(q)),
		}

		if id, err := xid.FromString(q); err == nil {
			predicates = append(predicates, account_ent.IDEQ(id))
		}

		aq.Where(account_ent.Or(predicates...))
	}
}

func WithSortBy(rule sortrule.SortRule) Filter {
	return func(aq *ent.AccountQuery) {
		if rule.Field() == "" {
			return
		}

		switch rule.Field() {
		case "created_at":
			if rule.IsDescending() {
				aq.Order(ent.Desc(account_ent.FieldCreatedAt))
			} else {
				aq.Order(ent.Asc(account_ent.FieldCreatedAt))
			}
		case "updated_at":
			if rule.IsDescending() {
				aq.Order(ent.Desc(account_ent.FieldUpdatedAt))
			} else {
				aq.Order(ent.Asc(account_ent.FieldUpdatedAt))
			}
		case "name":
			if rule.IsDescending() {
				aq.Order(ent.Desc(account_ent.FieldName))
			} else {
				aq.Order(ent.Asc(account_ent.FieldName))
			}
		case "handle":
			if rule.IsDescending() {
				aq.Order(ent.Desc(account_ent.FieldHandle))
			} else {
				aq.Order(ent.Asc(account_ent.FieldHandle))
			}
		}
	}
}

func WithRoles(roleIDs []xid.ID) Filter {
	return func(aq *ent.AccountQuery) {
		for _, roleID := range roleIDs {
			aq.Where(account_ent.HasAccountRolesWith(accountroles.RoleIDEQ(roleID)))
		}
	}
}

func WithJoinedInRange(tr timerange.TimeRange) Filter {
	return func(aq *ent.AccountQuery) {
		tr.Start.Call(func(start time.Time) {
			aq.Where(account_ent.CreatedAtGTE(start))
		})

		tr.End.Call(func(end time.Time) {
			aq.Where(account_ent.CreatedAtLTE(end))
		})
	}
}

func WithInvitedByHandles(handles []string) Filter {
	return func(aq *ent.AccountQuery) {
		if len(handles) == 0 {
			return
		}

		aq.Where(account_ent.HasInvitedByWith(
			invitation.HasCreatorWith(account_ent.HandleIn(handles...)),
		))
	}
}

func WithSuspended(suspended bool) Filter {
	return func(aq *ent.AccountQuery) {
		if suspended {
			aq.Where(account_ent.DeletedAtNotNil())
		} else {
			aq.Where(account_ent.DeletedAtIsNil())
		}
	}
}

func WithAdmin(admin bool) Filter {
	return func(aq *ent.AccountQuery) {
		aq.Where(account_ent.Admin(admin))
	}
}

func WithAuthServices(services []string) Filter {
	return func(aq *ent.AccountQuery) {
		if len(services) == 0 {
			return
		}

		aq.Where(account_ent.HasAuthenticationWith(auth_ent.ServiceIn(services...)))
	}
}

func (q *Querier) Search(ctx context.Context, params pagination.Parameters, filters ...Filter) (*pagination.Result[*account.AccountWithEdges], error) {
	aq := q.db.Account.Query().
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator()
		}).
		WithAuthentication()

	for _, fn := range filters {
		fn(aq)
	}

	total, err := aq.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	aq.Limit(params.Limit()).Offset(params.Offset())

	results, err := aq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(results)*2)
	for _, acc := range results {
		roleTargets = append(roleTargets, hydrationTargets(acc)...)
	}
	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := dt.MapErr(results, account.MapAccount)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, items)
	return &result, nil
}

func hydrationTargets(acc *ent.Account) []*ent.Account {
	targets := []*ent.Account{acc}

	if invitedBy := acc.Edges.InvitedBy; invitedBy != nil {
		creator, err := invitedBy.Edges.CreatorOrErr()
		if err == nil {
			targets = append(targets, creator)
		}
	}

	return targets
}
