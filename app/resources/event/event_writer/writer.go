package event_writer

import (
	"context"
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_querier"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/location"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_event "github.com/Southclaws/storyden/internal/ent/event"
)

type Writer struct {
	db      *ent.Client
	querier *event_querier.Querier
}

func New(db *ent.Client, querier *event_querier.Querier) *Writer {
	return &Writer{db: db, querier: querier}
}

type Option func(*ent.EventMutation)

func WithName(s string) Option {
	return func(m *ent.EventMutation) {
		m.SetName(s)
	}
}

func WithSlug(s string) Option {
	return func(m *ent.EventMutation) {
		m.SetSlug(s)
	}
}

func WithDescription(s string) Option {
	return func(m *ent.EventMutation) {
		m.SetDescription(s)
	}
}

func WithTimeRange(tr event_ref.TimeRange) Option {
	return func(m *ent.EventMutation) {
		m.SetStartTime(tr.Start)
		m.SetEndTime(tr.End)
	}
}

func WithImage(id asset.AssetID) Option {
	return func(m *ent.EventMutation) {
		m.SetPrimaryImageID(id)
	}
}

func WithParticipationPolicy(pp participation.Policy) Option {
	return func(m *ent.EventMutation) {
		m.SetParticipationPolicy(pp.String())
	}
}

func WithVisibility(vis visibility.Visibility) Option {
	return func(m *ent.EventMutation) {
		m.SetVisibility(ent_event.Visibility(vis.String()))
	}
}

func WithLocation(loc location.Location) Option {
	return func(m *ent.EventMutation) {
		switch l := loc.(type) {
		case *location.Physical:
			m.SetLocationName(l.Name)
			l.Address.Call(m.SetLocationAddress)
			l.Latitude.Call(m.SetLocationLatitude)
			l.Longitude.Call(m.SetLocationLongitude)
			l.URL.Call(func(v url.URL) { m.SetLocationURL(v.String()) })

		case *location.Virtual:
			m.SetLocationName(l.Name)
			l.URL.Call(func(v url.URL) { m.SetLocationURL(v.String()) })
		}
	}
}

func WithCapacity(c int) Option {
	return func(m *ent.EventMutation) {
		m.SetCapacity(c)
	}
}

func WithMetadata(v map[string]any) Option {
	return func(m *ent.EventMutation) {
		m.SetMetadata(v)
	}
}

func (w *Writer) Create(ctx context.Context,
	name string,
	slug string,
	startTime time.Time,
	endTime time.Time,
	policy participation.Policy,
	vis visibility.Visibility,
	threadID post.ID,
	opts ...Option,
) (*event.Event, error) {
	create := w.db.Event.Create()
	mutation := create.Mutation()

	mutation.SetName(name)
	mutation.SetSlug(slug)
	mutation.SetStartTime(startTime)
	mutation.SetEndTime(endTime)
	mutation.SetParticipationPolicy(policy.String())
	mutation.SetVisibility(ent_event.Visibility(vis.String()))
	mutation.SetThreadID(xid.ID(threadID))

	for _, opt := range opts {
		opt(mutation)
	}

	evt, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, event_ref.NewID(evt.ID))
}

func (w *Writer) Update(ctx context.Context, mk event_ref.QueryKey, opts ...Option) (*event.Event, error) {
	update := w.db.Event.Update()

	update.Where(mk.Predicate())

	mutation := update.Mutation()

	for _, opt := range opts {
		opt(mutation)
	}

	_, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, mk)
}

func (w *Writer) Delete(ctx context.Context, mk event_ref.QueryKey) (*event.Event, error) {
	update := w.db.Event.Update().
		Where(mk.Predicate()).
		SetDeletedAt(time.Now())

	_, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, mk)
}
