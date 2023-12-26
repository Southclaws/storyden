import { Portal } from "@ark-ui/react";
import format from "date-fns/format";

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

import { styled } from "@/styled-system/jsx";

import { Props, useFeedItemMenu } from "./useFeedItemMenu";

export function FeedItemMenu(props: Props) {
  const { onCopyLink, shareEnabled, onShare, deleteEnabled, onDelete } =
    useFeedItemMenu(props);

  return (
    <Menu size="sm">
      <MenuTrigger asChild>
        <MoreAction />
      </MenuTrigger>
      <Portal>
        <MenuPositioner>
          <MenuContent lazyMount minW="36">
            <MenuItemGroup id="user">
              <MenuItemGroupLabel
                htmlFor="user"
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>{`Post by ${props.thread.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.thread.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </MenuItemGroupLabel>

              <MenuSeparator />

              <MenuItem id="copy-link" onClick={onCopyLink}>
                Copy link
              </MenuItem>

              {shareEnabled && (
                <MenuItem id="share" onClick={onShare}>
                  Share
                </MenuItem>
              )}

              {deleteEnabled && (
                <MenuItem id="delete" onClick={onDelete}>
                  Delete
                </MenuItem>
              )}
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
