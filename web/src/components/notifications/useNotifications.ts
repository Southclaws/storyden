import filter from "lodash/fp/filter";
import flow from "lodash/fp/flow";
import map from "lodash/fp/map";

import { handle } from "@/api/client";
import {
  notificationUpdate,
  useNotificationList,
  useNotificationUpdateMany,
} from "@/api/openapi-client/notifications";
import {
  Notification,
  NotificationListResult,
  NotificationStatus,
} from "@/api/openapi-schema";
import { formatNotification } from "@/lib/notification";
import { useSession } from "@/auth";

import { NotificationItem } from "./item";

export type Props = {
  initialData?: NotificationListResult;
  status: NotificationStatus;
};

export function useNotifications(props: Props) {
  const session = useSession();
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

  const markAllReadMutation = useNotificationUpdateMany();

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const unreads = filterUnread(data.notifications).length;

  const notifications = processNotifications(data.notifications);

  async function handleMarkAs(id: string, status: NotificationStatus) {
    await handle(async () => {
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

  async function handleMarkAllAsRead() {
    await handle(async () => {
      if (!data) {
        return;
      }

      const updated = data.notifications.map((n) => ({
        id: n.id,
        status: NotificationStatus.read,
      }));

      await markAllReadMutation.trigger({ notifications: updated });

      if (data) {
        const newList = {
          ...data,
          notifications: data.notifications.map((n) => ({
            ...n,
            status: "read" as NotificationStatus,
          })),
        };

        await mutate(newList);
      }
    });
  }

  async function handleOpen(id: string) {
    if (!session?.meta.notifications.markReadOnClick) {
      return;
    }

    await handleMarkAs(id, "read");
  }

  return {
    ready: true as const,
    data: {
      unreads,
      notifications,
    },
    handlers: {
      handleOpen,
      handleMarkAs,
      handleMarkAllAsRead,
    },
  };
}

const filterStatus = (s: NotificationStatus) =>
  filter<Notification>((n) => n.status === s);

const filterUnread = filterStatus("unread");

const mapToItems = map(mapToItem);

function mapToItem(n: Notification): NotificationItem {
  const content = formatNotification(n);
  const createdAt = new Date(n.created_at);
  const title = content.title;
  const isRead = n.status === "read";

  return {
    id: n.id,
    createdAt,
    title,
    description: content.description,
    body: content.body,
    url: content.url,
    isRead,
    source: n.source,
    item: n.item,
  };
}
