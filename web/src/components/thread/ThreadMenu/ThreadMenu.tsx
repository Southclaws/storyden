"use client";

import { Portal } from "@ark-ui/react";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import { format } from "date-fns/format";
import { parseAsBoolean, useQueryState } from "nuqs";

import { Thread } from "src/api/openapi-schema";
import { MoreAction } from "src/components/site/Action/More";

import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";

export type Props = {
  thread: Thread;
};

export function usePostMenu(props: Props) {
  const [editing, setEditing] = useQueryState("edit", parseAsBoolean);

  function handleCopyLink() {
    //
  }
  function handleShare() {
    //
  }

  function handleEdit() {
    setEditing(true);
  }
  function handleDelete() {
    //
  }

  return {
    isShareEnabled: true,
    isEditingEnabled: true,
    isDeletingEnabled: true,
    handlers: {
      handleCopyLink,
      handleShare,
      handleEdit,
      handleDelete,
    },
  };
}

export function ThreadMenu(props: Props) {
  const { isShareEnabled, isEditingEnabled, isDeletingEnabled, handlers } =
    usePostMenu(props);

  const { thread } = props;

  return (
    <Menu.Root lazyMount>
      <Menu.Trigger asChild>
        <MoreAction variant="subtle" size="xs" />
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
                <styled.span>{`Post by ${thread.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(thread.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={handlers.handleCopyLink}>
                <HStack gap="1">
                  <LinkIcon width="1.4em" /> Copy link
                </HStack>
              </Menu.Item>

              {isShareEnabled && (
                <Menu.Item value="share" onClick={handlers.handleShare}>
                  <HStack gap="1">
                    <ShareIcon width="1.4em" /> Share
                  </HStack>
                </Menu.Item>
              )}

              {isEditingEnabled && (
                <Menu.Item value="edit" onClick={handlers.handleEdit}>
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
