"use client";

import { Portal } from "@ark-ui/react";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import { CancelAction } from "@/components/site/Action/Cancel";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

import { Props, useThreadMenu } from "./useThreadMenu";

export function ThreadMenu(props: Props) {
  const {
    isSharingEnabled,
    isEditingEnabled,
    isDeletingEnabled,
    isConfirmingDelete,
    handlers,
  } = useThreadMenu(props);

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
                  {format(new Date(thread.createdAt), "yyyy-MM-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={handlers.handleCopyLink}>
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
                <Menu.Item value="edit" onClick={handlers.handleEdit}>
                  <HStack gap="1">
                    <PencilIcon width="1.4em" /> Edit
                  </HStack>
                </Menu.Item>
              )}

              {isDeletingEnabled &&
                (isConfirmingDelete ? (
                  <HStack gap="0">
                    <Menu.Item
                      className={menuItemColorPalette({ colorPalette: "red" })}
                      value="delete"
                      w="full"
                      closeOnSelect={false}
                    >
                      Are you sure?
                    </Menu.Item>

                    <Menu.Item
                      value="delete-cancel"
                      closeOnSelect={false}
                      asChild
                    >
                      <CancelAction
                        borderRadius="md"
                        onClick={handlers.handleCancelDelete}
                      />
                    </Menu.Item>
                  </HStack>
                ) : (
                  <Menu.Item
                    className={menuItemColorPalette({ colorPalette: "red" })}
                    value="delete"
                    closeOnSelect={false}
                    onClick={handlers.handleConfirmDelete}
                  >
                    <HStack gap="1">
                      <TrashIcon width="1.4em" /> Delete
                    </HStack>
                  </Menu.Item>
                ))}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
