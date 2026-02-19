package event

import (
	"time"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

var _ datagraph.Item = (*Event)(nil)

type Event struct {
	event_ref.Event

	Thread thread.Thread
}

func (e *Event) GetID() xid.ID                 { return xid.ID(e.ID) }
func (e *Event) GetKind() datagraph.Kind       { return datagraph.KindEvent }
func (e *Event) GetName() string               { return e.Name }
func (e *Event) GetSlug() string               { return e.Slug }
func (e *Event) GetContent() datagraph.Content { return e.Thread.Content }
func (e *Event) GetDesc() string               { return e.Description.Or(e.Thread.GetDesc()) }
func (e *Event) GetProps() map[string]any      { return e.Meta }
func (e *Event) GetAssets() []*asset.Asset {
	if a, ok := e.Asset.Get(); ok {
		return []*asset.Asset{&a}
	}
	return nil
}
func (e *Event) GetCreated() time.Time { return e.CreatedAt }
func (e *Event) GetUpdated() time.Time { return e.UpdatedAt }

func Map(in *ent.Event, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Event, error) {
	threadEdge, err := in.Edges.ThreadOrErr()
	if err != nil {
		return nil, err
	}

	thr, err := thread.Map(threadEdge, roleHydratorFn)
	if err != nil {
		return nil, err
	}

	evt, err := event_ref.Map(in, roleHydratorFn)
	if err != nil {
		return nil, err
	}

	return &Event{
		Event:  *evt,
		Thread: *thr,
	}, nil
}
