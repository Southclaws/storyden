import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { keyBy } from "lodash";

import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { MoreIcon } from "@/components/ui/icons/More";
import * as Menu from "@/components/ui/menu";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../Context";

type BlockItem = {
  type: LibraryPageBlock["type"];
  label: string;
};

const blocks: BlockItem[] = [
  { type: "title", label: "Title" },
  { type: "cover", label: "Cover" },
  { type: "content", label: "Content" },
  { type: "properties", label: "Properties" },
  { type: "tags", label: "Tags" },
  { type: "assets", label: "Assets" },
  { type: "link", label: "Link" },
  { type: "table", label: "Table" },
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
                  {block.label}
                </Menu.Item>
              );
            })}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
