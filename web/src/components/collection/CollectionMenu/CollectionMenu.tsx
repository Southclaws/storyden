"use client";

import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import { DeleteWithConfirmationMenuItem } from "@/components/site/DeleteConfirmationMenuItem";
import { LinkIcon } from "@/components/ui/icons/Link";
import { ShareIcon } from "@/components/ui/icons/Share";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";

import { Props, useCollectionMenu } from "./useCollectionMenu";

export function CollectionMenu(props: Props) {
  const { isSharingEnabled, isDeletingEnabled, handlers } =
    useCollectionMenu(props);

  const { collection } = props;

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
                <styled.span>{`Collection by ${collection.owner.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(collection.createdAt), "yyyy-MM-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={handlers.handleCopyLink}>
                <HStack gap="1">
                  <LinkIcon /> Copy link
                </HStack>
              </Menu.Item>

              {isSharingEnabled && (
                <Menu.Item value="share" onClick={handlers.handleShare}>
                  <HStack gap="1">
                    <ShareIcon /> Share
                  </HStack>
                </Menu.Item>
              )}

              {isDeletingEnabled && (
                <DeleteWithConfirmationMenuItem
                  onDelete={handlers.handleDelete}
                />
              )}

              {/* {isDeletingEnabled &&
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
                ))} */}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
