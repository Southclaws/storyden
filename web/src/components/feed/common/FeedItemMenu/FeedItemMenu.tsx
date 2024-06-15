import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import * as Menu from "@/components/ui/menu";
import { styled } from "@/styled-system/jsx";

import { Props, useFeedItemMenu } from "./useFeedItemMenu";

export function FeedItemMenu(props: Props) {
  const { shareEnabled, deleteEnabled, handleSelect } = useFeedItemMenu(props);

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <MoreAction />
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup id="user">
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>{`Post by ${props.thread.author.name}`}</styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.thread.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link">Copy link</Menu.Item>
              {shareEnabled && <Menu.Item value="share">Share</Menu.Item>}
              {deleteEnabled && <Menu.Item value="delete">Delete</Menu.Item>}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
