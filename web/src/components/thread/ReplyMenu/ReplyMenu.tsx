"use client";

import { Portal } from "@ark-ui/react";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";

import { Props, useReplyMenu } from "./useReplyMenu";

export function ReplyMenu(props: Props) {
  const { isSharingEnabled, isEditingEnabled, isDeletingEnabled, handlers } =
    useReplyMenu(props);

  return (
    <Menu.Root lazyMount>
      <Menu.Trigger asChild>
        <MoreAction size="xs" />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>{`Post by ${props.reply.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.reply.createdAt), "yyyy-MM-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={handlers.handleCopyURL}>
                <HStack gap="1">
                  <LinkIcon width="1.4em" /> Copy link
                </HStack>
              </Menu.Item>

              {isSharingEnabled && (
                <Menu.Item value="share" onClick={handlers.handleShare}>
                  <HStack gap="1">
                    <ShareIcon width="1.4em" /> Share
                  </HStack>
                </Menu.Item>
              )}

              {isEditingEnabled && (
                <Menu.Item value="edit" onClick={handlers.handleSetEditing}>
                  <HStack gap="1">
                    <PencilIcon width="1.4em" /> Edit
                  </HStack>
                </Menu.Item>
              )}

              {isDeletingEnabled && (
                <Menu.Item value="delete" onClick={handlers.handleDelete}>
                  <HStack gap="1">
                    <TrashIcon width="1.4em" /> Delete
                  </HStack>
                </Menu.Item>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
