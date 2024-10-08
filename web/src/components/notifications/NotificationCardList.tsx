import {
  ArchiveBoxIcon,
  InboxArrowDownIcon,
} from "@heroicons/react/24/outline";
import Image from "next/image";

import { handle } from "@/api/client";
import { NotificationStatus } from "@/api/openapi-schema";
import { Card, CardRows } from "@/components/ui/rich-card";
import { css } from "@/styled-system/css";
import { Center, HStack, LStack, styled } from "@/styled-system/jsx";
import { timestamp } from "@/utils/date";

import { MemberBadge } from "../member/MemberBadge/MemberBadge";
import { Button } from "../ui/button";

import { NotificationItem } from "./item";

type Props = {
  notifications: NotificationItem[];
  onMove: (id: string, status: NotificationStatus) => Promise<void>;
};

export function NotificationCardList({ notifications, onMove }: Props) {
  if (notifications.length === 0) {
    return (
      <Center h="96" w="full" display="flex" flexDirection="column" gap="1">
        <styled.p color="fg.muted">no notifications.</styled.p>
      </Center>
    );
  }

  return (
    <CardRows>
      {notifications.map((n) => {
        const title = n.item?.description
          ? `${n.description} "${n.item?.description}"`
          : n.description;

        return (
          <Card
            key={n.id}
            id={n.id}
            shape="row"
            title={timestamp(n.createdAt, false)}
            text={title}
            url={n.url}
            controls={<StatusControl notification={n} onMove={onMove} />}
          >
            <NotificationSource {...n} />
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
        <styled.span color="fg.subtle">system message</styled.span>
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
    <Button
      variant="ghost"
      size="sm"
      title="Mark as unread"
      onClick={handleChangeStatus}
    >
      <InboxArrowDownIcon />
    </Button>
  ) : (
    <Button
      variant="ghost"
      size="sm"
      title="Mark as read"
      onClick={handleChangeStatus}
    >
      <ArchiveBoxIcon />
    </Button>
  );
}
