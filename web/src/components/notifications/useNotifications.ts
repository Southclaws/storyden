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
  switch (n.event) {
    case "thread_reply":
      return { description: "replied to your post", url: `/t/${n.item?.slug}` };
    case "post_like":
      return { description: "liked your post", url: `/t/${n.item?.slug}` };
    case "profile_mention":
      return { description: "mentioned you", url: `/t/${n.item?.slug}` };
    case "follow":
      return { description: "followed you", url: `/m/${n.source?.handle}` };
  }
}
