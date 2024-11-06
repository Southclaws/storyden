import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { CancelAction } from "@/components/site/Action/Cancel";
import { MoreAction } from "@/components/site/Action/More";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LinkIcon } from "@/components/ui/icons/Link";
import { ShareIcon } from "@/components/ui/icons/Share";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

import { Props, useFeedItemMenu } from "./useFeedItemMenu";

export function FeedItemMenu(props: Props) {
  const { isSharingEnabled, isDeletingEnabled, isConfirmingDelete, handlers } =
    useFeedItemMenu(props);

  return (
    <Menu.Root lazyMount onSelect={handlers.handleSelect}>
      <Menu.Trigger asChild>
        <MoreAction variant="subtle" size="xs" />
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
                  {format(new Date(props.thread.createdAt), "yyyy-MM-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link">
                <HStack gap="1">
                  <LinkIcon /> Copy link
                </HStack>
              </Menu.Item>

              {isSharingEnabled && (
                <Menu.Item value="share">
                  <HStack gap="1">
                    <ShareIcon /> Share
                  </HStack>
                </Menu.Item>
              )}

              {isDeletingEnabled &&
                (isConfirmingDelete ? (
                  <HStack gap="0">
                    <Menu.Item
                      className={menuItemColorPalette({ colorPalette: "red" })}
                      value="delete"
                      w="full"
                      closeOnSelect={false}
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
