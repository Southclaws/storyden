package notification

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Notification struct {
	ID     xid.ID
	Event  Event
	Item   opt.Optional[Item]
	Source opt.Optional[Source]
	Time   time.Time
	Read   bool
	Target opt.Optional[string]
}

type Item struct {
	ID    xid.ID
	Kind  datagraph.Kind
	Title string
	Slug  opt.Optional[string]
}

type Source struct {
	ID     account.AccountID
	Handle string
	Name   string
}

type Notifications []*Notification

func (a Notifications) Len() int           { return len(a) }
func (a Notifications) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Notifications) Less(i, j int) bool { return a[i].Time.After(a[j].Time) }

type NotificationRef struct {
	ID      xid.ID
	Event   Event
	ItemRef opt.Optional[datagraph.Ref]
	Source  opt.Optional[Source]
	Time    time.Time
	Read    bool
	Target  opt.Optional[string]
}

type NotificationRefs []*NotificationRef

func Map(r *ent.Notification) (*NotificationRef, error) {
	sourceEdge := opt.NewPtr(r.Edges.Source)

	source := opt.Map(sourceEdge, func(a ent.Account) Source {
		return Source{ID: account.AccountID(a.ID), Handle: a.Handle, Name: a.Name}
	})

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
		Target:  opt.NewPtr(r.Target),
	}, nil
}
