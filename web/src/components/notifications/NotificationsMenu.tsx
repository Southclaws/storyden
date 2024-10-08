import { Portal } from "@ark-ui/react";
import { ArchiveBoxIcon, Cog6ToothIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

import * as Menu from "@/components/ui/menu";
import { Box, Center, HStack, LStack, styled } from "@/styled-system/jsx";
import { hstack } from "@/styled-system/patterns";

import { MemberAvatar } from "../member/MemberBadge/MemberAvatar";
import { NotificationAction } from "../site/Navigation/Actions/Notifications";
import { Unready } from "../site/Unready";
import { Button } from "../ui/button";
import { LinkButton } from "../ui/link-button";

import { NotificationItem } from "./item";
import { Props, useNotifications } from "./useNotifications";

export function NotificationsMenu(props: Props) {
  const { ready, error, data, handlers } = useNotifications(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { unreads, notifications } = data;

  const isEmpty = notifications.length === 0;

  return (
    <Menu.Root closeOnSelect={false}>
      <Menu.Trigger cursor="pointer" position="relative">
        <NotificationAction hideLabel size="md" variant="ghost" />

        {!isEmpty && (
          <Box
            position="absolute"
            top="1"
            right="1"
            bgColor="red.8"
            borderRadius="full"
            w="2"
            h="2"
          />
        )}
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="heading">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <LStack fontSize="sm">
                  <HStack w="full" justify="space-between">
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
                  </HStack>
                </LStack>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              {isEmpty ? (
                <Center w="full" py="4" color="fg.muted" fontSize="xs">
                  You&apos;re all caught up!
                </Center>
              ) : (
                notifications.map((notification) => (
                  <Menu.Item value={notification.id} height="auto" py="1">
                    <HStack w="full" justify="space-between">
                      <Link
                        className={hstack({
                          w: "full",
                          justify: "space-between",
                        })}
                        href={notification.url}
                        onClick={() =>
                          handlers.handleMarkAs(notification.id, "read")
                        }
                      >
                        <NotificationAvatar notification={notification} />
                        <LStack gap="0">
                          <styled.span fontWeight="bold">
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
                        <ArchiveBoxIcon />
                      </Button>
                    </HStack>
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

  return <Cog6ToothIcon width="1rem" />;
}
