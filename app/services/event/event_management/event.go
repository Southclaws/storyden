package event_management

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_querier"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/event_writer"
	"github.com/Southclaws/storyden/app/resources/event/location"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/event/participation/participant_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var errNotAuthorised = fault.New("not authorised")

type Manager struct {
	accountQuery *account_querier.Querier
	querier      *event_querier.Querier
	writer       *event_writer.Writer
	partWriter   *participant_writer.Writer

	threadWriter thread.Service
	bus          *pubsub.Bus
}

func New(
	accountQuery *account_querier.Querier,
	querier *event_querier.Querier,
	writer *event_writer.Writer,
	partWriter *participant_writer.Writer,
	threadWriter thread.Service,
	bus *pubsub.Bus,
) *Manager {
	return &Manager{
		accountQuery: accountQuery,
		querier:      querier,
		writer:       writer,
		partWriter:   partWriter,
		threadWriter: threadWriter,
		bus:          bus,
	}
}

type Partial struct {
	Name                opt.Optional[string]
	Slug                opt.Optional[string]
	Description         opt.Optional[string]
	TimeRange           opt.Optional[event_ref.TimeRange]
	Image               opt.Optional[asset.AssetID]
	ParticipationPolicy opt.Optional[participation.Policy]
	Visibility          opt.Optional[visibility.Visibility]
	Location            opt.Optional[location.Location]
	Capacity            opt.Optional[int]
	Metadata            opt.Optional[map[string]any]
}

func (m *Manager) Create(ctx context.Context,
	name string,
	content datagraph.Content,
	startTime time.Time,
	endTime time.Time,
	policy participation.Policy,
	vis visibility.Visibility,
	cat category.CategoryID,
	partial Partial,
) (*event.Event, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []event_writer.Option{}

	partial.Description.Call(func(v string) { opts = append(opts, event_writer.WithDescription(v)) })
	partial.Location.Call(func(v location.Location) { opts = append(opts, event_writer.WithLocation(v)) })
	partial.Capacity.Call(func(v int) { opts = append(opts, event_writer.WithCapacity(v)) })
	partial.Metadata.Call(func(v map[string]any) { opts = append(opts, event_writer.WithMetadata(v)) })

	slug := partial.Slug.Or(mark.Slugify(name))

	thread, err := m.threadWriter.Create(
		ctx,
		name,
		accountID,
		nil,
		thread.Partial{
			Content:    opt.New(content),
			Category:   opt.New(xid.ID(cat)),
			Visibility: opt.New(visibility.VisibilityUnlisted),
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err := m.writer.Create(ctx, name, slug, startTime, endTime, policy, vis, thread.ID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create event"))
	}

	mk := event_ref.NewID(xid.ID(evt.ID))

	err = m.partWriter.Add(ctx,
		mk,
		accountID,
		participant_writer.WithRole(participation.RoleHost),
		participant_writer.WithStatus(participation.StatusAttending))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err = m.querier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventActivityCreated{
		ID: evt.ID,
	})

	if vis == visibility.VisibilityPublished {
		m.bus.Publish(ctx, &rpc.EventActivityPublished{
			ID: evt.ID,
		})
	}

	return evt, nil
}

func (m *Manager) Update(ctx context.Context, mk event_ref.QueryKey, partial Partial) (*event.Event, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	current, err := m.querier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if !current.Participants.IsHost(accountID) {
			return fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []event_writer.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, event_writer.WithName(v)) })
	partial.Slug.Call(func(v string) { opts = append(opts, event_writer.WithSlug(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, event_writer.WithDescription(v)) })
	partial.TimeRange.Call(func(v event_ref.TimeRange) { opts = append(opts, event_writer.WithTimeRange(v)) })
	partial.Image.Call(func(v asset.AssetID) { opts = append(opts, event_writer.WithImage(v)) })
	partial.ParticipationPolicy.Call(func(v participation.Policy) { opts = append(opts, event_writer.WithParticipationPolicy(v)) })
	partial.Visibility.Call(func(v visibility.Visibility) { opts = append(opts, event_writer.WithVisibility(v)) })
	partial.Location.Call(func(v location.Location) { opts = append(opts, event_writer.WithLocation(v)) })
	partial.Capacity.Call(func(v int) { opts = append(opts, event_writer.WithCapacity(v)) })
	partial.Metadata.Call(func(v map[string]any) { opts = append(opts, event_writer.WithMetadata(v)) })

	evt, err := m.writer.Update(ctx, mk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to update event"))
	}

	m.bus.Publish(ctx, &rpc.EventActivityUpdated{
		ID: evt.ID,
	})

	if vis, ok := partial.Visibility.Get(); ok && vis == visibility.VisibilityPublished {
		if current.Visibility != evt.Visibility {
			m.bus.Publish(ctx, &rpc.EventActivityPublished{
				ID: evt.ID,
			})
		}
	}

	return evt, nil
}

func (m *Manager) Delete(ctx context.Context, mk event_ref.QueryKey, partial Partial) (*event.Event, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	current, err := m.querier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if !current.Participants.IsHost(accountID) {
			return fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	evt, err := m.writer.Delete(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventActivityDeleted{
		ID: evt.ID,
	})

	return evt, nil
}
