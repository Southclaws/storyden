import { Portal } from "@ark-ui/react";
import { PositioningOptions } from "@zag-js/popper";
import { keyBy } from "lodash";

import { AddIcon } from "@/components/ui/icons/Add";
import * as Menu from "@/components/ui/menu";
import { BlockIcon } from "@/lib/library/blockIcons";
import { allBlockTypes } from "@/lib/library/blockTypes";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockName } from "@/lib/library/metadata";

import { useWatch } from "../store";

export function CreateBlockMenu({
  trigger,
  positioning = undefined,
  index = undefined,
}: {
  trigger?: React.ReactElement;
  positioning?: PositioningOptions;
  index?: number;
}) {
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = useWatch((s) => s.draft.meta);

  function handleAddBlock(type: LibraryPageBlock["type"]) {
    emit("library:add-block", {
      type,
      index: index ?? undefined,
    });
  }

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = allBlockTypes.filter((b) => !existingBlocks[b]);

  return (
    <Menu.Root lazyMount positioning={positioning}>
      {trigger ? (
        <Menu.Trigger asChild>{trigger}</Menu.Trigger>
      ) : (
        <Menu.TriggerItem>
          <AddIcon />
          &nbsp;Add Block
        </Menu.TriggerItem>
      )}

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            {blockList.map((block) => {
              return (
                <Menu.Item
                  key={block}
                  value={block}
                  onClick={() => handleAddBlock(block)}
                >
                  <BlockIcon blockType={block} />
                  &nbsp;
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
