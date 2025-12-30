package profile_search

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/sortrule"
	"github.com/Southclaws/storyden/app/resources/timerange"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/invitation"
)

type Filter func(*ent.AccountQuery)

type Querier struct {
	db *ent.Client
}

func WithDisplayNameContains(q string) Filter {
	return func(pq *ent.AccountQuery) {
		pq.Where(account.NameContainsFold(q))
	}
}

func WithHandleContains(q string) Filter {
	return func(pq *ent.AccountQuery) {
		pq.Where(account.HandleContainsFold(q))
	}
}

func WithNamesLike(q string) Filter {
	return func(pq *ent.AccountQuery) {
		pq.Where(account.Or(
			account.NameContainsFold(q),
			account.HandleContainsFold(q),
		))
	}
}

func WithSortBy(rule sortrule.SortRule) Filter {
	return func(pq *ent.AccountQuery) {
		if rule.Field() == "" {
			return
		}

		switch rule.Field() {
		case "created_at":
			if rule.IsDescending() {
				pq.Order(ent.Desc(account.FieldCreatedAt))
			} else {
				pq.Order(ent.Asc(account.FieldCreatedAt))
			}
		case "name":
			if rule.IsDescending() {
				pq.Order(ent.Desc(account.FieldName))
			} else {
				pq.Order(ent.Asc(account.FieldName))
			}
		case "handle":
			if rule.IsDescending() {
				pq.Order(ent.Desc(account.FieldHandle))
			} else {
				pq.Order(ent.Asc(account.FieldHandle))
			}
		}
	}
}

func WithRoles(roleIDs []xid.ID) Filter {
	return func(pq *ent.AccountQuery) {
		if len(roleIDs) == 0 {
			return
		}

		for _, roleID := range roleIDs {
			pq.Where(account.HasAccountRolesWith(
				accountroles.RoleIDEQ(roleID),
			))
		}
	}
}

func WithJoinedInRange(tr timerange.TimeRange) Filter {
	return func(pq *ent.AccountQuery) {
		tr.Start.Call(func(start time.Time) {
			pq.Where(account.CreatedAtGTE(start))
		})

		tr.End.Call(func(end time.Time) {
			pq.Where(account.CreatedAtLTE(end))
		})
	}
}

func WithInvitedBy(accountIDs []xid.ID) Filter {
	return func(pq *ent.AccountQuery) {
		if len(accountIDs) == 0 {
			return
		}

		pq.Where(account.HasInvitedByWith(
			invitation.CreatorAccountIDIn(accountIDs...),
		))
	}
}

func WithInvitedByHandles(handles []string) Filter {
	return func(pq *ent.AccountQuery) {
		if len(handles) == 0 {
			return
		}

		pq.Where(account.HasInvitedByWith(
			invitation.HasCreatorWith(
				account.HandleIn(handles...),
			),
		))
	}
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (d *Querier) Search(ctx context.Context, params pagination.Parameters, filters ...Filter) (*pagination.Result[*profile.Public], error) {
	q := d.db.Account.Query().
		WithTags().
		WithEmails().
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator()
		}).
		WithAuthentication()

	for _, fn := range filters {
		fn(q)
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	q.Limit(params.Limit()).Offset(params.Offset())

	r, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	profiles, err := dt.MapErr(r, profile.Map(nil))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, profiles)

	return &result, nil
}
