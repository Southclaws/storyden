"use client";

import { Portal } from "@ark-ui/react";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import { format } from "date-fns/format";

import { Post } from "src/api/openapi/schemas";
import { MoreAction } from "src/components/site/Action/More";

import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";

import { usePostMenu } from "./usePostMenu";

export function PostMenu(props: Post) {
  const {
    onCopyLink,
    shareEnabled,
    onShare,
    editEnabled,
    onEdit,
    deleteEnabled,
    onDelete,
  } = usePostMenu(props);

  return (
    <Menu.Root lazyMount>
      <Menu.Trigger asChild>
        <MoreAction />
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
                <styled.span>{`Post by ${props.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={onCopyLink}>
                <HStack gap="1">
                  <LinkIcon width="1.4em" /> Copy link
                </HStack>
              </Menu.Item>

              {shareEnabled && (
                <Menu.Item value="share" onClick={onShare}>
                  <HStack gap="1">
                    <ShareIcon width="1.4em" /> Share
                  </HStack>
                </Menu.Item>
              )}

              {/* <Menu.Item>Reply</Menu.Item> */}

              {editEnabled && (
                <Menu.Item value="edit" onClick={onEdit}>
                  <HStack gap="1">
                    <PencilIcon width="1.4em" /> Edit
                  </HStack>
                </Menu.Item>
              )}

              {deleteEnabled && (
                <Menu.Item value="delete" onClick={onDelete}>
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
