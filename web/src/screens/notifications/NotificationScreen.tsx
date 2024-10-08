"use client";

import { useQueryState } from "nuqs";

import {
  NotificationListResult,
  NotificationStatus,
} from "@/api/openapi-schema";
import { NotificationCardList } from "@/components/notifications/NotificationCardList";
import { useNotifications } from "@/components/notifications/useNotifications";
import { UnreadyBanner } from "@/components/site/Unready";
import { Switch } from "@/components/ui/switch";
import { HStack, LStack, styled } from "@/styled-system/jsx";

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

  return (
    <LStack>
      <HStack w="full" justify="space-between">
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
      </HStack>

      <NotificationCardList
        notifications={notifications}
        onMove={handlers.handleMarkAs}
      />
    </LStack>
  );
}
