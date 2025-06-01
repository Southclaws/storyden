import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { keyBy } from "lodash";

import { AddIcon } from "@/components/ui/icons/Add";
import * as Menu from "@/components/ui/menu";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../Context";

type BlockItem = {
  type: LibraryPageBlock["type"];
};

const blocks: BlockItem[] = [
  { type: "title" },
  { type: "cover" },
  { type: "content" },
  { type: "properties" },
  { type: "tags" },
  { type: "assets" },
  { type: "link" },
  { type: "table" },
];

export function CreateBlockMenu() {
  const { node, form } = useLibraryPageContext();
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = form.watch("meta", node.meta);

  function handleSelect(value: MenuSelectionDetails) {
    emit("library:add-block", {
      type: value.value as LibraryPageBlock["type"],
    });
  }

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = blocks.filter((b) => !existingBlocks[b.type]);

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
                <Menu.Item key={block.type} value={block.type}>
                  {LibraryPageBlockName[block.type]}
                </Menu.Item>
              );
            })}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
