import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { keyBy } from "lodash";

import { AddIcon } from "@/components/ui/icons/Add";
import * as Menu from "@/components/ui/menu";
import { allBlockTypes } from "@/lib/library/blockTypes";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";

import { useWatch } from "../store";

export function CreateBlockMenu() {
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = useWatch((s) => s.draft.meta);

  function handleSelect(value: MenuSelectionDetails) {
    emit("library:add-block", {
      type: value.value as LibraryPageBlock["type"],
    });
  }

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = allBlockTypes.filter((b) => !existingBlocks[b]);

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="add">
          <AddIcon />
          &nbsp;Add Block
        </Menu.Item>
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            {blockList.map((block) => {
              return (
                <Menu.Item key={block} value={block}>
                  {LibraryPageBlockName[block]}
                </Menu.Item>
              );
            })}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
