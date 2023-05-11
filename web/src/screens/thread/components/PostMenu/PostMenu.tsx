import { PostProps } from "src/api/openapi/schemas";
import { usePostMenu } from "./usePostMenu";
import {
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuDivider,
  MenuGroup,
} from "@chakra-ui/react";
import { ShareIcon } from "@heroicons/react/24/solid";
import { More } from "src/components/Action/Action";
import { LinkIcon, PencilIcon, TrashIcon } from "@heroicons/react/24/outline";
import format from "date-fns/format";

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
      <MenuButton as={More} />
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
