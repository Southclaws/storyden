import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { keyBy } from "lodash";
import { PropsWithChildren } from "react";

import { ButtonProps } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import * as Menu from "@/components/ui/menu";
import { allBlockTypes } from "@/lib/library/blockTypes";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";
import { styled } from "@/styled-system/jsx";

import { useWatch } from "../store";

import { CreateBlockMenu } from "./CreateBlockMenu";
import { LibraryPageAssetsBlockMenuItems } from "./LibraryPageAssetsBlock/LibraryPageAssetsBlockMenuItems";
import { LibraryPageCoverBlockMenuItems } from "./LibraryPageCoverBlock/LibraryPageCoverBlockMenuItems";
import { LibraryPageDirectoryBlockMenuItems } from "./LibraryPageDirectoryBlock/LibraryPageDirectoryBlockMenuItems";
import { LibraryPageTitleBlockMenuItems } from "./LibraryPageTitleBlock/LibraryPageTitleBlockMenuItems";

type Props = {
  open?: boolean;
  block: LibraryPageBlock;
  index: number;
};

type AllProps = PropsWithChildren<Props & ButtonProps>;

export function BlockMenu({ children, open, block, index }: AllProps) {
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = useWatch((s) => s.draft.meta);

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = allBlockTypes.filter((b) => !existingBlocks[b]);

  const newBlocksAvailable = blockList.length > 0;

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
              <BlockConfigMenu index={index} block={block} />
              {newBlocksAvailable && <CreateBlockMenu />}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function BlockConfigMenu({ block }: Props) {
  switch (block.type) {
    case "assets": {
      return <LibraryPageAssetsBlockMenuItems />;
    }
    case "cover": {
      return <LibraryPageCoverBlockMenuItems />;
    }
    case "directory": {
      return <LibraryPageDirectoryBlockMenuItems />;
    }
    case "title": {
      return <LibraryPageTitleBlockMenuItems />;
    }
    default:
      return null;
  }
}
