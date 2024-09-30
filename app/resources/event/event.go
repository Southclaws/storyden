package event

import (
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
func (e *Event) GetKind() datagraph.Kind       { return datagraph.KindPost }
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

func Map(in *ent.Event) (*Event, error) {
	threadEdge, err := in.Edges.ThreadOrErr()
	if err != nil {
		return nil, err
	}

	thr, err := thread.FromModel(nil, nil)(threadEdge)
	if err != nil {
		return nil, err
	}

	evt, err := event_ref.Map(in)
	if err != nil {
		return nil, err
	}

	return &Event{
		Event:  *evt,
		Thread: *thr,
	}, nil
}
