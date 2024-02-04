import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

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
  const { shareEnabled, deleteEnabled, handleSelect } = useFeedItemMenu(props);

  return (
    <Menu size="sm" lazyMount onSelect={handleSelect}>
      <MenuTrigger asChild>
        <MoreAction />
      </MenuTrigger>
      <Portal>
        <MenuPositioner>
          <MenuContent minW="36">
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

              <MenuItem id="copy-link">Copy link</MenuItem>
              {shareEnabled && <MenuItem id="share">Share</MenuItem>}
              {deleteEnabled && <MenuItem id="delete">Delete</MenuItem>}
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
