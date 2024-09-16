package notify_querier

import (
	"context"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entnotification "github.com/Southclaws/storyden/internal/ent/notification"
)

type Querier struct {
	db           *ent.Client
	postSearcher post_search.Repository
}

func New(db *ent.Client, postSearcher post_search.Repository) *Querier {
	return &Querier{db: db, postSearcher: postSearcher}
}

func (n *Querier) ListNotifications(ctx context.Context, accountID account.AccountID) (notification.Notifications, error) {
	r, err := n.db.Notification.Query().
		Where(entnotification.HasOwnerWith(entaccount.ID(xid.ID(accountID)))).
		WithSource(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	refs, err := dt.MapErr(r, notification.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ns, err := n.hydrateRefs(ctx, refs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ns, nil
}

func (n *Querier) hydrateRefs(ctx context.Context, refs notification.NotificationRefs) (notification.Notifications, error) {
	grouped := lo.GroupBy(refs, func(n *notification.NotificationRef) datagraph.Kind {
		return n.ItemRef.OrZero().Kind
	})

	pids := dt.Map(grouped[datagraph.KindPost], func(n *notification.NotificationRef) post.ID {
		return post.ID(n.ItemRef.OrZero().ID)
	})
	posts, err := n.postSearcher.GetMany(ctx, pids...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	pg := lo.KeyBy(posts, func(p *post.Post) post.ID { return p.ID })

	ns := dt.Map(refs, func(r *notification.NotificationRef) *notification.Notification {
		switch r.ItemRef.OrZero().Kind {
		case datagraph.KindPost:
			p := pg[post.ID(r.ItemRef.OrZero().ID)]
			return &notification.Notification{
				ID:     r.ID,
				Event:  r.Event,
				Item:   p,
				Source: r.Source,
				Time:   r.Time,
				Read:   r.Read,
			}
		}

		return &notification.Notification{
			ID:     r.ID,
			Event:  r.Event,
			Source: r.Source,
			Time:   r.Time,
			Read:   r.Read,
		}
	})

	sort.Sort(notification.Notifications(ns))

	return ns, nil
}
