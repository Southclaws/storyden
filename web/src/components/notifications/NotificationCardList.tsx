import { handle } from "@/api/client";
import { NotificationStatus } from "@/api/openapi-schema";
import { ArchiveIcon } from "@/components/ui/icons/Archive";
import { InboxIcon } from "@/components/ui/icons/Inbox";
import { Card, CardRows } from "@/components/ui/rich-card";
import { Text } from "@/components/ui/text";
import { getCommonProperties } from "@/lib/datagraph/item";
import { Center, HStack, LStack, WStack } from "@/styled-system/jsx";
import { timestamp } from "@/utils/date";

import { MemberBadge } from "../member/MemberBadge/MemberBadge";
import { Button } from "../ui/button";
import { IconButton } from "../ui/icon-button";

import { NotificationItem } from "./item";

type Props = {
  notifications: NotificationItem[];
  onMove: (id: string, status: NotificationStatus) => Promise<void>;
};

export function NotificationCardList({ notifications, onMove }: Props) {
  if (notifications.length === 0) {
    return (
      <Center h="96" w="full" display="flex" flexDirection="column" gap="1">
        <Text variant="supporting">no notifications.</Text>
      </Center>
    );
  }

  return (
    <CardRows>
      {notifications.map((n) => {
        const properties = n.item && getCommonProperties(n.item);

        const title = properties?.description
          ? `${n.description} "${properties?.description}"`
          : n.description;

        return (
          <Card
            key={n.id}
            id={n.id}
            shape="row"
            title={timestamp(n.createdAt, false)}
            text={title}
            url={n.url}
            // controls={}
          >
            <WStack>
              <NotificationSource {...n} />
              <StatusControl notification={n} onMove={onMove} />
            </WStack>
          </Card>
        );
      })}
    </CardRows>
  );
}

function NotificationSource(props: NotificationItem) {
  if (props.source) {
    return (
      <MemberBadge profile={props.source} size="sm" name="full-horizontal" />
    );
  }

  return (
    <HStack>
      <LStack gap="0">
        <Text as="span" variant="metadata">
          system message
        </Text>
      </LStack>
    </HStack>
  );
}

function StatusControl({
  notification,
  onMove,
}: {
  notification: NotificationItem;
  onMove: (id: string, status: NotificationStatus) => void;
}) {
  function handleChangeStatus() {
    handle(async () => {
      const newStatus = notification.isRead ? "unread" : "read";
      onMove(notification.id, newStatus);
    });
  }

  return notification.isRead ? (
    <IconButton
      variant="ghost"
      title="Mark as unread"
      onClick={handleChangeStatus}
    >
      <InboxIcon color="text.muted" />
    </IconButton>
  ) : (
    <IconButton
      variant="ghost"
      title="Mark as read"
      onClick={handleChangeStatus}
    >
      <ArchiveIcon color="text.muted" />
    </IconButton>
  );
}
