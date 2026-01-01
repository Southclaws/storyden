"use client";

import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import { DatagraphItemKind } from "@/api/openapi-schema";
import { CategoryMoveMenu } from "@/components/category/CategoryMoveMenu/CategoryMoveMenu";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import {
  ReportPostMenuItem,
  truncateBody,
} from "@/components/report/ReportPostMenuItem";
import { CancelAction } from "@/components/site/Action/Cancel";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { EditIcon } from "@/components/ui/icons/Edit";
import { LinkIcon } from "@/components/ui/icons/Link";
import { PinIcon, PinOffIcon } from "@/components/ui/icons/Pin";
import { ShareIcon } from "@/components/ui/icons/Share";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

import { Props, useThreadMenu } from "./useThreadMenu";

export function ThreadMenu(props: Props) {
  const {
    isSharingEnabled,
    isEditingEnabled,
    isMovingEnabled,
    isDeletingEnabled,
    isConfirmingDelete,
    canPinThread,
    isThreadPinned,
    handlers,
  } = useThreadMenu(props);

  const { thread } = props;

  return (
    <Menu.Root
      positioning={{
        shift: 32,
      }}
      lazyMount
    >
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

                <MemberBadge
                  profile={thread.author}
                  size="sm"
                  name="full-vertical"
                />

                <styled.time fontWeight="normal">
                  {format(new Date(thread.createdAt), "yyyy-MM-dd")}
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

              <ReportPostMenuItem
                menuLabel="Report thread"
                targetKind={DatagraphItemKind.thread}
                targetId={thread.id}
                author={thread.author}
                headline={thread.title || "Untitled thread"}
                body={truncateBody(thread.description)}
              />

              {canPinThread && !isThreadPinned && (
                <Menu.Item value="pin" onClick={handlers.handlePinThread}>
                  <HStack gap="1">
                    <PinIcon /> Pin thread
                  </HStack>
                </Menu.Item>
              )}

              {canPinThread && isThreadPinned && (
                <Menu.Item value="unpin" onClick={handlers.handleUnpinThread}>
                  <HStack gap="1">
                    <PinOffIcon /> Unpin thread
                  </HStack>
                </Menu.Item>
              )}

              {isEditingEnabled && (
                <Menu.Item value="edit" onClick={handlers.handleEdit}>
                  <HStack gap="1">
                    <EditIcon /> Edit
                  </HStack>
                </Menu.Item>
              )}

              {isMovingEnabled && <CategoryMoveMenu thread={thread} />}

              {isDeletingEnabled &&
                (isConfirmingDelete ? (
                  <HStack gap="0">
                    <Menu.Item
                      className={menuItemColorPalette({ colorPalette: "red" })}
                      value="delete"
                      w="full"
                      closeOnSelect={false}
                      onClick={handlers.handleConfirmDelete}
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
                      <DeleteIcon /> Delete
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
