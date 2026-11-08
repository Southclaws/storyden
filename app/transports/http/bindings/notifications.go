package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/notification/push"
	"github.com/Southclaws/storyden/app/services/notification/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
)

type Notifications struct {
	notifyReader *notify_querier.Querier
	notifyWriter *notify_writer.Writer
	push         *push.Service
	pushRepo     *push_subscription.Repository
	settings     *settings.Service
}

func NewNotifications(
	notifyReader *notify_querier.Querier,
	notifyWriter *notify_writer.Writer,
	push *push.Service,
	pushRepo *push_subscription.Repository,
	settings *settings.Service,
) Notifications {
	return Notifications{
		notifyReader: notifyReader,
		notifyWriter: notifyWriter,
		push:         push,
		pushRepo:     pushRepo,
		settings:     settings,
	}
}

func (h *Notifications) NotificationList(ctx context.Context, request openapi.NotificationListRequestObject) (openapi.NotificationListResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := h.notifyReader.ListNotifications(ctx, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	notifications := dt.Map(list, serialiseNotification)

	return openapi.NotificationList200JSONResponse{
		NotificationListOKJSONResponse: openapi.NotificationListOKJSONResponse{
			// TODO: Pagination at resource layer.
			CurrentPage:   1,
			NextPage:      nil,
			Notifications: notifications,
			PageSize:      len(notifications),
			Results:       len(notifications),
			TotalPages:    1,
		},
	}, nil
}

func (h *Notifications) NotificationUpdate(ctx context.Context, request openapi.NotificationUpdateRequestObject) (openapi.NotificationUpdateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if request.Body.Status == nil {
		return nil, nil
	}

	id := openapi.ParseID(request.NotificationId)

	status := opt.Map(opt.NewPtr(request.Body.Status), deserialiseNotificationStatus)

	n, err := h.notifyWriter.SetRead(ctx, accountID, id, status.OrZero())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NotificationUpdate200JSONResponse{
		NotificationUpdateOKJSONResponse: openapi.NotificationUpdateOKJSONResponse(serialiseNotificationRef(n)),
	}, nil
}

func (h *Notifications) NotificationUpdateMany(ctx context.Context, request openapi.NotificationUpdateManyRequestObject) (openapi.NotificationUpdateManyResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	notificationRefs := make([]*notification.NotificationRef, 0, len(request.Body.Notifications))
	for _, n := range request.Body.Notifications {
		if n.Status == nil {
			continue
		}

		notificationRefs = append(notificationRefs, &notification.NotificationRef{
			ID:   openapi.ParseID(n.Id),
			Read: deserialiseNotificationStatus(*n.Status),
		})
	}

	updated, err := h.notifyWriter.UpdateStatusMany(ctx, accountID, notificationRefs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	notifications := dt.Map(updated, serialiseNotificationRef)

	return openapi.NotificationUpdateMany200JSONResponse{
		NotificationUpdateManyOKJSONResponse: openapi.NotificationUpdateManyOKJSONResponse{
			CurrentPage:   1,
			NextPage:      nil,
			Notifications: notifications,
			PageSize:      len(notifications),
			Results:       len(notifications),
			TotalPages:    1,
		},
	}, nil
}

func (h *Notifications) NotificationSettingsGet(ctx context.Context, request openapi.NotificationSettingsGetRequestObject) (openapi.NotificationSettingsGetResponseObject, error) {
	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	settings, err := h.settings.Get(ctx, owner)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NotificationSettingsGet200JSONResponse{
		NotificationSettingsGetOKJSONResponse: openapi.NotificationSettingsGetOKJSONResponse(serialiseNotificationSettings(settings)),
	}, nil
}

func (h *Notifications) NotificationSettingsUpdate(ctx context.Context, request openapi.NotificationSettingsUpdateRequestObject) (openapi.NotificationSettingsUpdateResponseObject, error) {
	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	settings, err := h.settings.Update(ctx, owner, deserialiseNotificationSettings(*request.Body))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NotificationSettingsUpdate200JSONResponse{
		NotificationSettingsUpdateOKJSONResponse: openapi.NotificationSettingsUpdateOKJSONResponse(serialiseNotificationSettings(settings)),
	}, nil
}

func (h *Notifications) PushSubscriptionUpsert(
	ctx context.Context,
	request openapi.PushSubscriptionUpsertRequestObject,
) (openapi.PushSubscriptionUpsertResponseObject, error) {
	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	subscription, err := h.push.Register(ctx, owner, webpush.SubscriptionInput{
		Endpoint: request.Body.Endpoint,
		P256DH:   request.Body.Keys.P256dh,
		Auth:     request.Body.Keys.Auth,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PushSubscriptionUpsert200JSONResponse{
		NotificationPushSubscriptionUpsertOKJSONResponse: openapi.NotificationPushSubscriptionUpsertOKJSONResponse{
			Id:       subscription.ID.String(),
			Endpoint: subscription.Endpoint.String(),
		},
	}, nil
}

func (h *Notifications) PushSubscriptionDelete(
	ctx context.Context,
	request openapi.PushSubscriptionDeleteRequestObject,
) (openapi.PushSubscriptionDeleteResponseObject, error) {
	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := h.pushRepo.Delete(ctx, owner, openapi.ParseID(request.PushSubscriptionId)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PushSubscriptionDelete204Response{}, nil
}

func serialiseNotification(in *notification.Notification) openapi.Notification {
	item := opt.Map(opt.NewSafe(in.Item, in.Item != nil), serialiseDatagraphItem)

	return openapi.Notification{
		Id:        in.ID.String(),
		CreatedAt: in.Time,
		Event:     openapi.NotificationEvent(in.Event.String()),
		Item:      item.Ptr(),
		Source:    opt.Map(in.Source, serialiseProfileReference).Ptr(),
		Status:    serialiseNotificationStatus(in.Read),
		Target:    in.Target.Ptr(),
	}
}

func serialiseNotificationRef(in *notification.NotificationRef) openapi.Notification {
	return openapi.Notification{
		Id:        in.ID.String(),
		CreatedAt: in.Time,
		Event:     openapi.NotificationEvent(in.Event.String()),
		Source:    opt.Map(in.Source, serialiseProfileReference).Ptr(),
		Status:    serialiseNotificationStatus(in.Read),
		Target:    in.Target.Ptr(),
	}
}

func serialiseNotificationStatus(in bool) openapi.NotificationStatus {
	if in {
		return openapi.Read
	}
	return openapi.Unread
}

func deserialiseNotificationStatus(in openapi.NotificationStatus) bool {
	return in == openapi.Read
}

func serialiseNotificationSettings(in account.NotificationPreferences) openapi.NotificationSettings {
	events := make(map[string]openapi.NotificationEventSettings)
	for event, settings := range notification.ResolveEventSettings(in.Events) {
		events[event] = openapi.NotificationEventSettings{
			Inbox: settings.Inbox,
			Push:  settings.Push,
		}
	}

	return openapi.NotificationSettings{
		Events: events,
		Push: openapi.PushNotificationSettings{
			Enabled: in.Push.Enabled,
		},
	}
}

func deserialiseNotificationSettings(in openapi.NotificationSettings) account.NotificationPreferences {
	events := make(map[string]account.NotificationEventSettings)
	for event, settings := range in.Events {
		events[event] = account.NotificationEventSettings{
			Inbox: settings.Inbox,
			Push:  settings.Push,
		}
	}

	return account.NotificationPreferences{
		Events: events,
		Push: account.PushNotificationPreferences{
			Enabled: in.Push.Enabled,
		},
	}
}
