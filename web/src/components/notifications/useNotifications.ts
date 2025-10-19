import { filter, flow, map } from "lodash/fp";

import { handle } from "@/api/client";
import {
  notificationUpdate,
  useNotificationList,
} from "@/api/openapi-client/notifications";
import {
  Notification,
  NotificationListResult,
  NotificationStatus,
} from "@/api/openapi-schema";
import { getCommonProperties } from "@/lib/datagraph/item";

import { NotificationItem } from "./item";

export type Props = {
  initialData?: NotificationListResult;
  status: NotificationStatus;
};

export function useNotifications(props: Props) {
  const filterByStatus = filterStatus(props.status);
  const processNotifications = flow(filterByStatus, mapToItems);

  const { data, error, mutate } = useNotificationList(
    { status: [props.status] },
    {
      swr: {
        fallbackData: props.initialData,
        revalidateIfStale: true,
        revalidateOnReconnect: true,
      },
    },
  );
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const unreads = filterUnread(data.notifications).length;

  const notifications = processNotifications(data.notifications);

  async function handleMarkAs(id: string, status: NotificationStatus) {
    handle(async () => {
      await notificationUpdate(id, { status });

      if (data) {
        const newList = {
          ...data,
          notifications: data.notifications.map((n) => {
            if (n.id === id) {
              return { ...n, status };
            }
            return n;
          }),
        };

        await mutate(newList);
      }
    });
  }

  return {
    ready: true as const,
    data: {
      unreads,
      notifications,
    },
    handlers: {
      handleMarkAs,
    },
  };
}

const filterStatus = (s: NotificationStatus) =>
  filter<Notification>((n) => n.status === s);

const filterUnread = filterStatus("unread");

const mapToItems = map(mapToItem);

function mapToItem(n: Notification): NotificationItem {
  const content = getNotificationContent(n);
  const createdAt = new Date(n.created_at);
  const title = n.source?.handle ?? "System";
  const isRead = n.status === "read";

  return {
    id: n.id,
    createdAt,
    title,
    description: content.description,
    url: content.url,
    isRead,
    source: n.source,
    item: n.item,
  };
}

function getNotificationContent(n: Notification) {
  const p = n.item && getCommonProperties(n.item);
  switch (n.event) {
    case "thread_reply":
      return { description: "replied to your post", url: `/t/${p?.slug}` };
    case "post_like":
      return { description: "liked your post", url: `/t/${p?.slug}` };
    case "follow":
      return { description: "followed you", url: `/m/${n.source?.handle}` };
    case "profile_mention":
      return { description: "mentioned you", url: `/t/${p?.slug}` };
    case "event_host_added":
      return { description: "added you as an event host", url: `#` }; // not implemented
    case "member_attending_event":
      return { description: "is attending your event", url: `#` }; // not implemented
    case "member_declined_event":
      return { description: "declined your event", url: `#` }; // not implemented
    case "attendee_removed":
      return { description: "removed you from their event", url: `#` }; // not implemented
    case "report_submitted":
      return { description: "submitted a report", url: `/reports` };
    case "report_updated":
      return { description: "report status updated", url: `/reports` };
  }
}
