import {
  Menu,
  MenuButton,
  MenuDivider,
  MenuGroup,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";
import { LinkIcon, TrashIcon } from "@heroicons/react/24/outline";
import { ShareIcon } from "@heroicons/react/24/solid";
import format from "date-fns/format";

import { ThreadReference } from "src/api/openapi/schemas";
import { More } from "src/components/site/Action/Action";

import { useThreadMenu } from "./useThreadMenu";

export function ThreadMenu(props: ThreadReference) {
  const { onCopyLink, shareEnabled, onShare, deleteEnabled, onDelete } =
    useThreadMenu(props);

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

        {deleteEnabled && (
          <MenuItem icon={<TrashIcon width="1.4em" />} onClick={onDelete}>
            Delete
          </MenuItem>
        )}
      </MenuList>
    </Menu>
  );
}
