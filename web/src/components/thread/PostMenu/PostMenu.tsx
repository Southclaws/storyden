"use client";

import { Portal } from "@ark-ui/react";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import format from "date-fns/format";

import { PostProps } from "src/api/openapi/schemas";
import { MoreAction } from "src/components/site/Action/More";
import {
  Menu,
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuPositioner,
  MenuSeparator,
  MenuTrigger,
} from "src/theme/components/Menu";

import { HStack, styled } from "@/styled-system/jsx";

import { usePostMenu } from "./usePostMenu";

export function PostMenu(props: PostProps) {
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
    <Menu size="sm">
      <MenuTrigger>
        <MoreAction />
      </MenuTrigger>

      <Portal>
        <MenuPositioner>
          <MenuContent lazyMount minW="36">
            <MenuItemGroup id="group">
              <MenuItemGroupLabel
                htmlFor="user"
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>{`Post by ${props.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </MenuItemGroupLabel>

              <MenuSeparator />

              <MenuItem id="copy-link" onClick={onCopyLink}>
                <HStack gap="1">
                  <LinkIcon width="1.4em" /> Copy link
                </HStack>
              </MenuItem>

              {shareEnabled && (
                <MenuItem id="share" onClick={onShare}>
                  <HStack gap="1">
                    <ShareIcon width="1.4em" /> Share
                  </HStack>
                </MenuItem>
              )}

              {/* <MenuItem>Reply</MenuItem> */}

              {editEnabled && (
                <MenuItem id="edit" onClick={onEdit}>
                  <HStack gap="1">
                    <PencilIcon width="1.4em" /> Edit
                  </HStack>
                </MenuItem>
              )}

              {deleteEnabled && (
                <MenuItem id="delete" onClick={onDelete}>
                  <HStack gap="1">
                    <TrashIcon width="1.4em" /> Delete
                  </HStack>
                </MenuItem>
              )}
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
