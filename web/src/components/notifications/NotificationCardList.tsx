import { handle } from "@/api/client";
import { NotificationStatus } from "@/api/openapi-schema";
import { ArchiveIcon } from "@/components/ui/icons/Archive";
import { InboxIcon } from "@/components/ui/icons/Inbox";
import { Card, CardRows } from "@/components/ui/rich-card";
import { useI18n } from "@/i18n/provider";
import { getCommonProperties } from "@/lib/datagraph/item";
import { Center, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
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
  const { t } = useI18n();

  if (notifications.length === 0) {
    return (
      <Center h="96" w="full" display="flex" flexDirection="column" gap="1">
        <styled.p color="fg.muted">{t("no notifications.")}</styled.p>
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
  const { t } = useI18n();

  if (props.source) {
    return (
      <MemberBadge profile={props.source} size="sm" name="full-horizontal" />
    );
  }

  return (
    <HStack>
      <LStack gap="0">
        <styled.span color="fg.subtle">{t("system message")}</styled.span>
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
  const { t } = useI18n();

  function handleChangeStatus() {
    handle(async () => {
      const newStatus = notification.isRead ? "unread" : "read";
      onMove(notification.id, newStatus);
    });
  }

  return notification.isRead ? (
    <IconButton
      variant="ghost"
      size="xs"
      title={t("Mark as unread")}
      onClick={handleChangeStatus}
    >
      <InboxIcon color="fg.subtle" />
    </IconButton>
  ) : (
    <IconButton
      variant="ghost"
      size="xs"
      title={t("Mark as read")}
      onClick={handleChangeStatus}
    >
      <ArchiveIcon color="fg.subtle" />
    </IconButton>
  );
}
