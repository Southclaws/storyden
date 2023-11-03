"use client";

import { Portal } from "@ark-ui/react";

import { ThreadReference } from "src/api/openapi/schemas";
import { More } from "src/components/site/Action/Action";
import {
  Menu,
  MenuContent,
  MenuItem,
  MenuPositioner,
  MenuTrigger,
} from "src/theme/components/Menu";

// import { useFeedItemMenu } from "./useFeedItemMenu";

export function FeedItemMenu(props: ThreadReference) {
  // const { onCopyLink, shareEnabled, onShare, deleteEnabled, onDelete } =
  //   useFeedItemMenu(props);

  return (
    <Menu>
      <MenuTrigger asChild>
        <More />
      </MenuTrigger>
      <Portal>
        <MenuPositioner>
          <MenuContent>
            <MenuItem id="one">{props.id}</MenuItem>
            <MenuItem id="twi">{props.author.handle}</MenuItem>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
