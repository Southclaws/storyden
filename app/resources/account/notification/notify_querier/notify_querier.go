package notify_querier

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entnode "github.com/Southclaws/storyden/internal/ent/node"
	entnotification "github.com/Southclaws/storyden/internal/ent/notification"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (n *Querier) Get(ctx context.Context, owner account.AccountID, id xid.ID) (*notification.Notification, error) {
	query := n.db.Notification.Query().Where(
		entnotification.OwnerAccountIDEQ(xid.ID(owner)),
		entnotification.IDEQ(id),
	)

	items, err := n.query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fault.New("notification not found", ftag.With(ftag.NotFound))
	}

	return items[0], nil
}

func (n *Querier) ListNotifications(ctx context.Context, owner account.AccountID, statuses ...bool) (notification.Notifications, error) {
	query := n.db.Notification.Query().Where(entnotification.OwnerAccountIDEQ(xid.ID(owner)))

	if len(statuses) > 0 {
		filters := make([]predicate.Notification, 0, len(statuses))
		for _, read := range statuses {
			filters = append(filters, entnotification.ReadEQ(read))
		}

		query.Where(entnotification.Or(filters...))
	}

	query.Order(ent.Desc(entnotification.FieldCreatedAt))

	return n.query(ctx, query)
}

type row struct {
	ID            xid.ID    `sql:"id"`
	CreatedAt     time.Time `sql:"created_at"`
	EventType     string    `sql:"event_type"`
	Read          bool      `sql:"read"`
	Target        *string   `sql:"target"`
	DatagraphKind *string   `sql:"datagraph_kind"`
	DatagraphID   *xid.ID   `sql:"datagraph_id"`
	SourceID      *xid.ID   `sql:"source_id"`
	SourceHandle  *string   `sql:"source_handle"`
	SourceName    *string   `sql:"source_name"`
	PostID        *xid.ID   `sql:"post_id"`
	PostTitle     *string   `sql:"post_title"`
	RootTitle     *string   `sql:"root_title"`
	NodeID        *xid.ID   `sql:"node_id"`
	NodeName      *string   `sql:"node_name"`
	NodeSlug      *string   `sql:"node_slug"`
	ProfileID     *xid.ID   `sql:"profile_id"`
	ProfileName   *string   `sql:"profile_name"`
	ProfileHandle *string   `sql:"profile_handle"`
}

func (n *Querier) query(ctx context.Context, query *ent.NotificationQuery) (notification.Notifications, error) {
	var rows []row
	err := query.Modify(func(s *sql.Selector) {
		source := sql.Table(entaccount.Table).As("notification_source")
		post := sql.Table(entpost.Table).As("notification_post")
		root := sql.Table(entpost.Table).As("notification_root")
		node := sql.Table(entnode.Table).As("notification_node")
		profile := sql.Table(entaccount.Table).As("notification_profile")

		s.LeftJoin(source).On(s.C(entnotification.FieldSourceAccountID), source.C(entaccount.FieldID))
		s.LeftJoin(post).OnP(sql.And(
			sql.In(s.C(entnotification.FieldDatagraphKind), datagraph.KindPost.String(), datagraph.KindThread.String(), datagraph.KindReply.String()),
			sql.ColumnsEQ(s.C(entnotification.FieldDatagraphID), post.C(entpost.FieldID)),
		))
		s.LeftJoin(root).On(post.C(entpost.FieldRootPostID), root.C(entpost.FieldID))
		s.LeftJoin(node).OnP(sql.And(
			sql.EQ(s.C(entnotification.FieldDatagraphKind), datagraph.KindNode.String()),
			sql.ColumnsEQ(s.C(entnotification.FieldDatagraphID), node.C(entnode.FieldID)),
		))
		s.LeftJoin(profile).OnP(sql.And(
			sql.EQ(s.C(entnotification.FieldDatagraphKind), datagraph.KindProfile.String()),
			sql.ColumnsEQ(s.C(entnotification.FieldDatagraphID), profile.C(entaccount.FieldID)),
		))

		s.Select(
			s.C(entnotification.FieldID),
			s.C(entnotification.FieldCreatedAt),
			s.C(entnotification.FieldEventType),
			s.C(entnotification.FieldRead),
			s.C(entnotification.FieldTarget),
			s.C(entnotification.FieldDatagraphKind),
			s.C(entnotification.FieldDatagraphID),
		)
		s.AppendSelectAs(source.C(entaccount.FieldID), "source_id")
		s.AppendSelectAs(source.C(entaccount.FieldHandle), "source_handle")
		s.AppendSelectAs(source.C(entaccount.FieldName), "source_name")
		s.AppendSelectAs(post.C(entpost.FieldID), "post_id")
		s.AppendSelectAs(post.C(entpost.FieldTitle), "post_title")
		s.AppendSelectAs(root.C(entpost.FieldTitle), "root_title")
		s.AppendSelectAs(node.C(entnode.FieldID), "node_id")
		s.AppendSelectAs(node.C(entnode.FieldName), "node_name")
		s.AppendSelectAs(node.C(entnode.FieldSlug), "node_slug")
		s.AppendSelectAs(profile.C(entaccount.FieldID), "profile_id")
		s.AppendSelectAs(profile.C(entaccount.FieldName), "profile_name")
		s.AppendSelectAs(profile.C(entaccount.FieldHandle), "profile_handle")
	}).Scan(ctx, &rows)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := make(notification.Notifications, 0, len(rows))
	for _, row := range rows {
		item, err := row.notification()
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if item == nil {
			continue
		}

		items = append(items, item)
	}

	return items, nil
}

func (r row) notification() (*notification.Notification, error) {
	event, err := notification.NewEvent(r.EventType)
	if err != nil {
		return nil, err
	}

	result := &notification.Notification{
		ID:     r.ID,
		Event:  event,
		Time:   r.CreatedAt,
		Read:   r.Read,
		Target: opt.NewPtr(r.Target),
	}

	if r.SourceID != nil {
		result.Source = opt.New(notification.Source{
			ID:     account.AccountID(*r.SourceID),
			Handle: opt.NewPtr(r.SourceHandle).OrZero(),
			Name:   opt.NewPtr(r.SourceName).OrZero(),
		})
	}

	if r.DatagraphID == nil || r.DatagraphKind == nil {
		return result, nil
	}

	kind, err := datagraph.NewKind(*r.DatagraphKind)
	if err != nil {
		return nil, err
	}

	item := notification.Item{ID: *r.DatagraphID, Kind: kind}
	switch kind {
	case datagraph.KindPost, datagraph.KindThread, datagraph.KindReply:
		if r.PostID == nil {
			return nil, nil
		}

		item.Title = opt.NewPtr(r.PostTitle).OrZero()
		if r.RootTitle != nil {
			item.Title = *r.RootTitle
		}
	case datagraph.KindNode:
		if r.NodeID == nil {
			return nil, nil
		}

		item.Title = opt.NewPtr(r.NodeName).OrZero()
		item.Slug = opt.NewPtr(r.NodeSlug)
	case datagraph.KindProfile:
		if r.ProfileID == nil {
			return nil, nil
		}

		item.Title = opt.NewPtr(r.ProfileName).OrZero()
		item.Slug = opt.NewPtr(r.ProfileHandle)
	default:
		return result, nil
	}

	result.Item = opt.New(item)

	return result, nil
}
