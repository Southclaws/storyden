"use client";

import { Portal } from "@ark-ui/react";

import { CancelAction } from "@/components/site/Action/Cancel";
import { Badge } from "@/components/ui/badge";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { EditIcon } from "@/components/ui/icons/Edit";
import { WarningIcon } from "@/components/ui/icons/Warning";
import * as Menu from "@/components/ui/menu";
import * as Tooltip from "@/components/ui/tooltip";
import { HStack } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

type Props = {
  isModerator: boolean;
  postId: string;
  onAccept: (postId: string) => void;
  onEditAndAccept: (postId: string) => void;
  onDelete: (postId: string) => void;
  isConfirmingDelete: boolean;
  onCancelDelete: () => void;
};

export function PostReviewBadge({
  isModerator,
  postId,
  onAccept,
  onEditAndAccept,
  onDelete,
  isConfirmingDelete,
  onCancelDelete,
}: Props) {
  if (!isModerator) {
    return (
      <Tooltip.Root
        openDelay={0}
        positioning={{
          slide: true,
          shift: 16,
        }}
      >
        <Tooltip.Trigger asChild>
          <Badge
            variant="subtle"
            cursor="pointer"
            aria-label="Post is in review"
          >
            <WarningIcon />
            In review
          </Badge>
        </Tooltip.Trigger>
        <Portal>
          <Tooltip.Positioner>
            <Tooltip.Arrow>
              <Tooltip.ArrowTip />
            </Tooltip.Arrow>

            <Tooltip.Content p="2" borderRadius="2xl" maxW="xs">
              Your post has been flagged for review by a moderator. It will be
              visible to others once approved.
            </Tooltip.Content>
          </Tooltip.Positioner>
        </Portal>
      </Tooltip.Root>
    );
  }

  return (
    <Menu.Root
      positioning={{
        placement: "bottom-start",
      }}
      lazyMount
    >
      <Menu.Trigger asChild>
        <Badge
          variant="subtle"
          cursor="pointer"
          _hover={{
            borderColor: "colorPalette.6",
          }}
          aria-label="Post review actions"
        >
          <WarningIcon />
          In review
        </Badge>
      </Menu.Trigger>

      <Menu.Positioner>
        <Menu.Content minW="48">
          <Menu.ItemGroup id="review-actions">
            <Menu.Item
              value="accept"
              onClick={() => onAccept(postId)}
              aria-label="Accept post"
            >
              <HStack gap="1">
                <CheckIcon /> Accept
              </HStack>
            </Menu.Item>

            <Menu.Item
              value="edit-and-accept"
              onClick={() => onEditAndAccept(postId)}
              aria-label="Edit and accept post"
            >
              <HStack gap="1">
                <EditIcon /> Edit and Accept
              </HStack>
            </Menu.Item>

            <Menu.Separator />

            {isConfirmingDelete ? (
              <HStack gap="0">
                <Menu.Item
                  className={menuItemColorPalette({ colorPalette: "red" })}
                  value="confirm-delete"
                  w="full"
                  closeOnSelect={false}
                  onClick={() => onDelete(postId)}
                  aria-label="Confirm delete post"
                >
                  Are you sure?
                </Menu.Item>

                <Menu.Item
                  value="cancel-delete"
                  closeOnSelect={false}
                  asChild
                  aria-label="Cancel delete"
                >
                  <CancelAction borderRadius="md" onClick={onCancelDelete} />
                </Menu.Item>
              </HStack>
            ) : (
              <Menu.Item
                className={menuItemColorPalette({ colorPalette: "red" })}
                value="delete"
                closeOnSelect={false}
                onClick={() => onDelete(postId)}
                aria-label="Delete post"
              >
                <HStack gap="1">
                  <DeleteIcon /> Delete
                </HStack>
              </Menu.Item>
            )}
          </Menu.ItemGroup>
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  );
}
