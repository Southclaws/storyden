import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PositioningOptions } from "@zag-js/popper";

import { AddIcon } from "@/components/ui/icons/Add";
import * as Menu from "@/components/ui/menu";
import {
  AllFeedBlockTypes,
  FeedBlockName,
  FeedBlockType,
} from "@/lib/settings/feed";

import { useFeedBlockEditor } from "./Context";
import { FeedBlockIcon } from "./blockIcons";

type Props = {
  index?: number;
  positioning?: PositioningOptions;
  trigger?: React.ReactElement;
};

export function CreateBlockMenu({ index, positioning, trigger }: Props) {
  const { addBlock, feed } = useFeedBlockEditor();
  const existing = new Set(feed.blocks.map((block) => block.type));
  const available = AllFeedBlockTypes.filter((type) => !existing.has(type));

  function handleSelect({ value }: MenuSelectionDetails) {
    void addBlock(value as FeedBlockType, index);
  }

  if (available.length === 0) {
    return null;
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect} positioning={positioning}>
      {trigger ? (
        <Menu.Trigger asChild>{trigger}</Menu.Trigger>
      ) : (
        <Menu.TriggerItem>
          <AddIcon />
          Add block
        </Menu.TriggerItem>
      )}
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="40">
            {available.map((type) => (
              <Menu.Item key={type} value={type}>
                <FeedBlockIcon type={type} />
                {FeedBlockName[type]}
              </Menu.Item>
            ))}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
