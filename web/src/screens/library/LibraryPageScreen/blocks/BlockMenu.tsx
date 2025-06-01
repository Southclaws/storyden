import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { NodeWithChildren } from "@/api/openapi-schema";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { MoreIcon } from "@/components/ui/icons/More";
import * as Menu from "@/components/ui/menu";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, WithMetadata } from "@/lib/library/metadata";
import { styled } from "@/styled-system/jsx";

import { CreateBlockMenu } from "./CreateBlockMenu";

type Props = {
  node: WithMetadata<NodeWithChildren>;
  block: LibraryPageBlock;
};

export function BlockMenu({ node, block }: Props & ButtonProps) {
  const emit = useEmitLibraryBlockEvent();

  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "delete": {
        emit("library:remove-block", {
          type: block.type,
        });
      }
    }
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <IconButton variant="ghost" size="xs" width="5" height="5" padding="0">
          <MoreIcon width="3" />
        </IconButton>
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup>
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>{block.type}</styled.span>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="delete">
                <DeleteIcon />
                &nbsp;Delete
              </Menu.Item>
              <CreateBlockMenu />
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
