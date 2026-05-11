import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PositioningOptions } from "@zag-js/popper";
import { keyBy } from "lodash";

import { AddIcon } from "@/components/ui/icons/Add";
import * as Menu from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";
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
  const { t } = useI18n();
  const emit = useEmitLibraryBlockEvent();
  const menuTrigger = trigger ?? (
    <Menu.Item value="add">
      <AddIcon />
      &nbsp;{t("Add Block")}
    </Menu.Item>
  );

  const currentMetadata = useWatch((s) => s.draft.meta);

  function handleSelect(value: MenuSelectionDetails) {
    emit("library:add-block", {
      type: value.value as LibraryPageBlock["type"],
      index: index ?? undefined,
    });
  }

  const existingBlocks = keyBy(currentMetadata.layout?.blocks, (b) => b.type);

  const blockList = allBlockTypes.filter((b) => !existingBlocks[b]);

  return (
    <Menu.Root lazyMount onSelect={handleSelect} positioning={positioning}>
      <Menu.Trigger asChild>{menuTrigger}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            {blockList.map((block) => {
              return (
                <Menu.Item key={block} value={block}>
                  <BlockIcon blockType={block} />
                  &nbsp;
                  {t(LibraryPageBlockName[block])}
                </Menu.Item>
              );
            })}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
