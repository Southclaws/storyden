package bindings

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/event"
	"github.com/Southclaws/storyden/app/resources/event/event_querier"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/event/location"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/services/event/event_management"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Events struct {
	eventQuerier *event_querier.Querier
	eventManager *event_management.Manager
}

func NewEvents(
	eventQuerier *event_querier.Querier,
	eventManager *event_management.Manager,
) Events {
	return Events{
		eventQuerier: eventQuerier,
		eventManager: eventManager,
	}
}

func (h *Events) EventList(ctx context.Context, request openapi.EventListRequestObject) (openapi.EventListResponseObject, error) {
	events, err := h.eventQuerier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.EventList200JSONResponse{
		EventListOKJSONResponse: openapi.EventListOKJSONResponse{
			// TODO: Pagination
			CurrentPage: 1,
			Events:      dt.Map(events, serialiseEventRefPtr),
			PageSize:    len(events),
			Results:     len(events),
			TotalPages:  1,
		},
	}, nil
}

func (h *Events) EventCreate(ctx context.Context, request openapi.EventCreateRequestObject) (openapi.EventCreateResponseObject, error) {
	content, err := datagraph.NewRichText(request.Body.Content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	participationPolicy, err := participation.NewPolicy(string(request.Body.ParticipationPolicy))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	visibility, err := deserialiseVisibility(request.Body.Visibility)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	location, err := opt.MapErr(opt.NewPtr(request.Body.Location), deserialiseLocation)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	catID := category.CategoryID(openapi.ParseID(request.Body.ThreadCategoryId))

	partial := event_management.Partial{
		Slug:        opt.NewPtr(request.Body.Slug),
		Description: opt.NewPtr(request.Body.Description),
		Image:       opt.Map(opt.NewPtr(request.Body.PrimaryImageAssetId), deserialiseAssetID),
		Capacity:    opt.NewPtr(request.Body.Capacity),
		Location:    location,
	}

	evt, err := h.eventManager.Create(ctx,
		request.Body.Name,
		content,
		request.Body.TimeRange.Start,
		request.Body.TimeRange.End,
		participationPolicy,
		visibility,
		catID,
		partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.EventCreate200JSONResponse{
		EventCreateOKJSONResponse: openapi.EventCreateOKJSONResponse(serialiseEventPtr(evt)),
	}, nil
}

func (h *Events) EventDelete(ctx context.Context, request openapi.EventDeleteRequestObject) (openapi.EventDeleteResponseObject, error) {
	// event_manager
	return nil, nil
}

func (h *Events) EventGet(ctx context.Context, request openapi.EventGetRequestObject) (openapi.EventGetResponseObject, error) {
	mk := event_ref.NewKey(request.EventMark)

	evt, err := h.eventQuerier.Get(ctx, mk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.EventGet200JSONResponse{
		EventGetOKJSONResponse: openapi.EventGetOKJSONResponse(serialiseEventPtr(evt)),
	}, nil
}

func (h *Events) EventUpdate(ctx context.Context, request openapi.EventUpdateRequestObject) (openapi.EventUpdateResponseObject, error) {
	// event_manager
	return nil, nil
}

func (h *Events) EventParticipantRemove(ctx context.Context, request openapi.EventParticipantRemoveRequestObject) (openapi.EventParticipantRemoveResponseObject, error) {
	// participation_manager
	return nil, nil
}

func (h *Events) EventParticipantUpdate(ctx context.Context, request openapi.EventParticipantUpdateRequestObject) (openapi.EventParticipantUpdateResponseObject, error) {
	// participation_manager
	return nil, nil
}

func serialiseEventPtr(in *event.Event) openapi.Event {
	return openapi.Event{
		Id:                  in.ID.String(),
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
		DeletedAt:           in.DeletedAt.Ptr(),
		Name:                in.Name,
		Slug:                in.Slug,
		Description:         in.Description.String(),
		TimeRange:           serialiseEventTimeRange(in.TimeRange),
		PrimaryImage:        opt.Map(in.Asset, serialiseAsset).Ptr(),
		Participants:        serialiseParticipants(in.Participants),
		ParticipationPolicy: serialiseParticipationPolicy(in.Policy),
		Visibility:          serialiseVisibility(in.Visibility),
		Location:            serialiseLocation(in.Location),
		Capacity:            in.Capacity.Ptr(),
		Thread:              serialiseThread(&in.Thread),
		Meta:                (*openapi.Metadata)(&in.Meta),
	}
}

func serialiseEventRef(in event_ref.Event) openapi.EventReference {
	return openapi.EventReference{
		Id:                  in.ID.String(),
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
		DeletedAt:           in.DeletedAt.Ptr(),
		Name:                in.Name,
		Slug:                in.Slug,
		Description:         in.Description.String(),
		TimeRange:           serialiseEventTimeRange(in.TimeRange),
		PrimaryImage:        opt.Map(in.Asset, serialiseAsset).Ptr(),
		Participants:        serialiseParticipants(in.Participants),
		ParticipationPolicy: serialiseParticipationPolicy(in.Policy),
		Visibility:          serialiseVisibility(in.Visibility),
		Location:            serialiseLocation(in.Location),
		Capacity:            in.Capacity.Ptr(),
		Meta:                (*openapi.Metadata)(&in.Meta),
	}
}

func serialiseEventRefPtr(in *event_ref.Event) openapi.EventReference {
	return serialiseEventRef(*in)
}

func serialiseEventTimeRange(in event_ref.TimeRange) openapi.EventTimeRange {
	return openapi.EventTimeRange{
		Start: in.Start,
		End:   in.End,
	}
}

func serialiseParticipantPtr(in *participation.EventParticipant) openapi.EventParticipant {
	return openapi.EventParticipant{
		Profile: serialiseProfileReferencePtr(&in.Account),
		Role:    openapi.EventParticipantRole(in.Role.String()),
		Status:  openapi.EventParticipationStatus(in.Status.String()),
	}
}

func serialiseParticipants(in participation.EventParticipants) []openapi.EventParticipant {
	return dt.Map(in, serialiseParticipantPtr)
}

func serialiseParticipationPolicy(in participation.Policy) openapi.EventParticipationPolicy {
	return openapi.EventParticipationPolicy(in.String())
}

func deserialiseLocation(in openapi.EventLocation) (location.Location, error) {
	lt, err := in.ValueByDiscriminator()
	if err != nil {
		return nil, err
	}

	switch l := lt.(type) {
	case openapi.EventLocationPhysical:
		u, err := opt.MapErr(opt.NewPtr(l.Url), deserialiseURL)
		if err != nil {
			return nil, err
		}

		return &location.Physical{
			Name:      l.Name,
			Address:   opt.NewPtr(l.Address),
			Latitude:  deserialiseOptionalFloat(l.Latitude),
			Longitude: deserialiseOptionalFloat(l.Longitude),
			URL:       u,
		}, nil

	case openapi.EventLocationVirtual:
		u, err := opt.MapErr(opt.NewPtr(l.Url), deserialiseURL)
		if err != nil {
			return nil, err
		}

		return &location.Virtual{
			Name: l.Name,
			URL:  u,
		}, nil

	default:
		return nil, fault.Newf("invalid location type %T", lt)
	}
}

func deserialiseURL(in string) (url.URL, error) {
	u, err := url.Parse(in)
	if err != nil {
		return url.URL{}, err
	}

	return *u, nil
}

func serialiseLocation(in location.Location) openapi.EventLocation {
	l := openapi.EventLocation{}
	var err error
	switch loc := in.(type) {
	case *location.Physical:
		err = l.FromEventLocationPhysical(openapi.EventLocationPhysical{
			Name:      loc.Name,
			Address:   loc.Address.Ptr(),
			Latitude:  serialiseOptionalFloat(loc.Latitude),
			Longitude: serialiseOptionalFloat(loc.Longitude),
			Url:       seraliseOptionalURL(loc.URL),
		})

	case *location.Virtual:
		err = l.FromEventLocationVirtual(openapi.EventLocationVirtual{
			Name: loc.Name,
			Url:  seraliseOptionalURL(loc.URL),
		})

	default:
		// Do nothing, invalid but not end of the world.
	}

	if err != nil {
		slog.Error("failed to serialise event location", slog.String("error", err.Error()), slog.Any("location", in))
	}

	return l
}
