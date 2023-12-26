"use client";

import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import format from "date-fns/format";

import { PostProps } from "src/api/openapi/schemas";
import { MoreAction } from "src/components/site/Action/More";
import {
  Menu,
  MenuButton,
  MenuDivider,
  MenuGroup,
  MenuItem,
  MenuList,
} from "src/theme/components";

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
    <Menu>
      <MenuButton>
        <MoreAction />
      </MenuButton>
      <MenuList>
        <MenuGroup title={`Post by ${props.author.name}`}>
          <MenuItem isDisabled>
            {format(new Date(props.createdAt), "yyyy-mm-dd")}
          </MenuItem>
        </MenuGroup>
        <MenuDivider />

        <MenuItem icon={<LinkIcon width="1.4em" />} onClick={onCopyLink}>
          Copy link
        </MenuItem>

        {shareEnabled && (
          <MenuItem icon={<ShareIcon width="1.4em" />} onClick={onShare}>
            Share
          </MenuItem>
        )}

        {/* <MenuItem>Reply</MenuItem> */}

        {editEnabled && (
          <MenuItem icon={<PencilIcon width="1.4em" />} onClick={onEdit}>
            Edit
          </MenuItem>
        )}

        {deleteEnabled && (
          <MenuItem icon={<TrashIcon width="1.4em" />} onClick={onDelete}>
            Delete
          </MenuItem>
        )}
      </MenuList>
    </Menu>
  );
}
