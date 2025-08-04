import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { ButtonProps } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import * as Menu from "@/components/ui/menu";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";
import { styled } from "@/styled-system/jsx";

import { CreateBlockMenu } from "./CreateBlockMenu";

type Props = {
  open?: boolean;
  block: LibraryPageBlock;
  index: number;
};

type AllProps = PropsWithChildren<Props & ButtonProps>;

export function BlockMenu({ children, open, block, index }: AllProps) {
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
    <Menu.Root
      open={open}
      lazyMount
      onSelect={handleSelect}
      positioning={{
        placement: "right-start",
        gutter: 0,
      }}
    >
      <Menu.Trigger asChild>
        {/*  */}
        {children}
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
                <styled.span>{LibraryPageBlockName[block.type]}</styled.span>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="delete">
                <DeleteIcon />
                &nbsp;Delete
              </Menu.Item>
              <CreateBlockMenu index={index} />
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
