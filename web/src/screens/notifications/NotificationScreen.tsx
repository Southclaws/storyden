"use client";

import { useQueryState } from "nuqs";

import {
  NotificationListResult,
  NotificationStatus,
} from "@/api/openapi-schema";
import { NotificationCardList } from "@/components/notifications/NotificationCardList";
import { useNotifications } from "@/components/notifications/useNotifications";
import { UnreadyBanner } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  initialData: NotificationListResult;
};

export function useNotificationScreen(props: Props) {
  const [status, setStatus] = useQueryState<NotificationStatus>("status", {
    defaultValue: "unread",
    parse(v: string) {
      switch (v) {
        case "read":
          return NotificationStatus.read;
        default:
          return NotificationStatus.unread;
      }
    },
  });
  const { ready, data, error, handlers } = useNotifications({
    initialData: props.initialData,
    status,
  });
  if (!ready) {
    return {
      ready: false as const,
      error,
    };
  }

  function handleToggleStatus() {
    setStatus(
      status === NotificationStatus.unread
        ? NotificationStatus.read
        : NotificationStatus.unread,
    );
  }

  return {
    ready: true as const,
    data,
    status,
    handlers: {
      handleToggleStatus,
      handleMarkAs: handlers.handleMarkAs,
      handleMarkAllAsRead: handlers.handleMarkAllAsRead,
    },
  };
}

export function NotificationScreen(props: Props) {
  const { ready, error, data, status, handlers } = useNotificationScreen(props);
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const { notifications } = data;

  const showingArchived = status === NotificationStatus.read;

  const hasUnreadNotifications = data.unreads > 0 && !showingArchived;

  return (
    <LStack>
      <WStack justifyContent="space-between" alignItems="flex-start">
        <LStack>
          <styled.h1 fontWeight="bold">Notifications</styled.h1>

          <Switch
            size="sm"
            checked={showingArchived}
            onClick={handlers.handleToggleStatus}
          >
            Archived
          </Switch>
        </LStack>

        {hasUnreadNotifications && (
          <Button
            variant="ghost"
            size="sm"
            flexShrink="0"
            onClick={handlers.handleMarkAllAsRead}
          >
            Mark all as read
          </Button>
        )}
      </WStack>

      <NotificationCardList
        notifications={notifications}
        onMove={handlers.handleMarkAs}
      />
    </LStack>
  );
}
