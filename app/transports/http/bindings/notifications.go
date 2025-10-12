package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Notifications struct {
	notifyReader *notify_querier.Querier
	notifyWriter *notify_writer.Writer
}

func NewNotifications(
	notifyReader *notify_querier.Querier,
	notifyWriter *notify_writer.Writer,
) Notifications {
	return Notifications{
		notifyReader: notifyReader,
		notifyWriter: notifyWriter,
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
			NextPage:      new(int),
			Notifications: notifications,
			PageSize:      len(notifications),
			Results:       len(notifications),
			TotalPages:    1,
		},
	}, nil
}

func (h *Notifications) NotificationUpdate(ctx context.Context, request openapi.NotificationUpdateRequestObject) (openapi.NotificationUpdateResponseObject, error) {
	// NOTE: Ownership is not checked, only authentication, so if you know the
	// ID of someone elses notification (unlikely) you can mark it as read.
	_, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if request.Body.Status == nil {
		return nil, nil
	}

	id := openapi.ParseID(request.NotificationId)

	status := opt.Map(opt.NewPtr(request.Body.Status), deserialiseNotificationStatus)

	n, err := h.notifyWriter.SetRead(ctx, id, status.OrZero())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NotificationUpdate200JSONResponse{
		NotificationUpdateOKJSONResponse: openapi.NotificationUpdateOKJSONResponse(serialiseNotificationRef(n)),
	}, nil
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
	}
}

func serialiseNotificationRef(in *notification.NotificationRef) openapi.Notification {
	return openapi.Notification{
		Id:        in.ID.String(),
		CreatedAt: in.Time,
		Event:     openapi.NotificationEvent(in.Event.String()),
		Source:    opt.Map(in.Source, serialiseProfileReference).Ptr(),
		Status:    serialiseNotificationStatus(in.Read),
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
