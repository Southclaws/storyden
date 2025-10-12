import { Portal } from "@ark-ui/react";
import Link from "next/link";

import { MemberAvatar } from "@/components/member/MemberBadge/MemberAvatar";
import { Button } from "@/components/ui/button";
import { ArchiveIcon } from "@/components/ui/icons/Archive";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { LinkButton } from "@/components/ui/link-button";
import * as Menu from "@/components/ui/menu";
import { Center, LStack, WStack, styled } from "@/styled-system/jsx";
import { hstack } from "@/styled-system/patterns";
import { deriveError } from "@/utils/error";

import { NotificationsTrigger } from "./NotificationsTrigger";
import { NotificationItem } from "./item";
import { Props, useNotifications } from "./useNotifications";

export function NotificationsMenu(props: Props) {
  const { ready, error, data, handlers } = useNotifications(props);
  if (!ready) {
    return (
      <NotificationsTrigger
        hideLabel
        size="md"
        variant="ghost"
        disabled
        title={deriveError(error)}
      />
    );
  }

  const { unreads, notifications } = data;

  const isEmpty = notifications.length === 0;

  return (
    <Menu.Root closeOnSelect={false}>
      <Menu.Trigger cursor="pointer" position="relative" asChild>
        <NotificationsTrigger
          hideLabel
          size="md"
          variant="ghost"
          unread={!isEmpty}
        />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="heading">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <LStack fontSize="sm">
                  <WStack>
                    <styled.p color="fg.muted">
                      Notifications ({unreads})
                    </styled.p>

                    <LinkButton
                      href="/notifications"
                      size="xs"
                      variant="outline"
                    >
                      see all
                    </LinkButton>
                  </WStack>
                </LStack>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              {isEmpty ? (
                <Center w="full" py="4" color="fg.muted" fontSize="xs">
                  You&apos;re all caught up!
                </Center>
              ) : (
                notifications.map((notification) => (
                  <Menu.Item
                    key={notification.id}
                    value={notification.id}
                    height="auto"
                    py="1"
                  >
                    <WStack className="notification-menu__row" minW="0" gap="1">
                      <Link
                        className={hstack({
                          w: "full",
                          minW: "0",
                          justify: "space-between",
                        })}
                        href={notification.url}
                        onClick={() =>
                          handlers.handleMarkAs(notification.id, "read")
                        }
                      >
                        <NotificationAvatar notification={notification} />
                        <LStack gap="0" minW="0">
                          <styled.span
                            fontWeight="bold"
                            textWrap="nowrap"
                            textOverflow="ellipsis"
                            overflow="hidden"
                            maxW="full"
                          >
                            {notification.source?.handle ?? "System"}
                          </styled.span>
                          <styled.span fontWeight="normal">
                            {notification.description}
                          </styled.span>
                        </LStack>
                      </Link>

                      <Button
                        variant="ghost"
                        size="sm"
                        title="Mark as read"
                        onClick={() =>
                          handlers.handleMarkAs(notification.id, "read")
                        }
                      >
                        <ArchiveIcon />
                      </Button>
                    </WStack>
                  </Menu.Item>
                ))
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

export function NotificationAvatar(props: { notification: NotificationItem }) {
  if (props.notification.source) {
    return <MemberAvatar profile={props.notification.source} size="sm" />;
  }

  return <SettingsIcon />;
}
