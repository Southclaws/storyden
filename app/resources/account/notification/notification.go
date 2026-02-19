package notification

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type Notification struct {
	ID     xid.ID
	Event  Event
	Item   datagraph.Item
	Source opt.Optional[profile.Ref]
	Time   time.Time
	Read   bool
}

type Notifications []*Notification

func (a Notifications) Len() int           { return len(a) }
func (a Notifications) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Notifications) Less(i, j int) bool { return a[i].Time.After(a[j].Time) }

type NotificationRef struct {
	ID      xid.ID
	Event   Event
	ItemRef opt.Optional[datagraph.Ref]
	Source  opt.Optional[profile.Ref]
	Time    time.Time
	Read    bool
}

type NotificationRefs []*NotificationRef

func Map(roleHydratorFn func(accID xid.ID) (held.Roles, error)) func(r *ent.Notification) (*NotificationRef, error) {
	return func(r *ent.Notification) (*NotificationRef, error) {
		sourceEdge := opt.NewPtr(r.Edges.Source)

		source, err := opt.MapErr(sourceEdge, func(a ent.Account) (profile.Ref, error) {
			p, err := profile.RefMapper(roleHydratorFn)(&a)
			if err != nil {
				return profile.Ref{}, err
			}
			return *p, err
		})
		if err != nil {
			return nil, err
		}

		et, err := NewEvent(r.EventType)
		if err != nil {
			return nil, err
		}

		var itemRef opt.Optional[datagraph.Ref]
		if r.DatagraphKind != nil && r.DatagraphID != nil {
			k, err := datagraph.NewKind(*r.DatagraphKind)
			if err != nil {
				return nil, err
			}

			itemRef = opt.New(datagraph.Ref{ID: *r.DatagraphID, Kind: k})
		}

		return &NotificationRef{
			ID:      r.ID,
			Event:   et,
			ItemRef: itemRef,
			Source:  source,
			Time:    r.CreatedAt,
			Read:    r.Read,
		}, nil
	}
}
