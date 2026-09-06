import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ReactNode, useCallback, useState } from "react";

import * as BlockEditor from "@/components/ui/block-editor";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DragItemFeedBlock } from "@/lib/dragdrop/provider";
import { useFeedBlockEvent } from "@/lib/feed/events";
import { FeedBlock, FeedBlockType } from "@/lib/settings/feed";
import { Box, LStack } from "@/styled-system/jsx";

import { BlockMenu } from "./BlockMenu";
import { useFeedBlockEditor } from "./Context";
import { CreateBlockMenu } from "./CreateBlockMenu";
import { FeedCategoriesBlock } from "./blocks/FeedCategoriesBlock";
import { FeedContentBlock } from "./blocks/FeedContentBlock";
import { FeedCoverBlock } from "./blocks/FeedCoverBlock";
import { FeedLibraryBlock } from "./blocks/FeedLibraryBlock";
import { FeedQuickShareBlock } from "./blocks/FeedQuickShareBlock";
import { FeedSubtitleBlock } from "./blocks/FeedSubtitleBlock";
import { FeedThreadsBlock } from "./blocks/FeedThreadsBlock";
import { FeedTitleBlock } from "./blocks/FeedTitleBlock";

export function FeedPageBlocks() {
  const { feed, isEditing, moveBlock } = useFeedBlockEditor();

  const handleReorder = useCallback(
    (activeId: FeedBlockType, overId: FeedBlockType) => {
      void moveBlock(activeId, overId);
    },
    [moveBlock],
  );
  useFeedBlockEvent("feed:reorder-block", ({ activeId, overId }) => {
    handleReorder(activeId, overId);
  });

  const blockIDs = feed.blocks.map((block) =>
    getFeedBlockSortableID(block.type),
  );

  if (!isEditing) {
    return (
      <LStack className="feed-page__blocks" gap="4" width="full">
        {feed.blocks.map((block) => (
          <FeedBlockRender key={block.type} block={block} />
        ))}
      </LStack>
    );
  }

  return (
    <LStack className="feed-page__blocks" gap="4" width="full">
      <SortableContext items={blockIDs} strategy={verticalListSortingStrategy}>
        {feed.blocks.map((block, index) => (
          <FeedBlockEditable key={block.type} block={block} index={index} />
        ))}
      </SortableContext>

      <CreateBlockMenu
        trigger={
          <Button variant="outline" width="full">
            <AddIcon />
            Add block
          </Button>
        }
        positioning={{ placement: "bottom" }}
      />
    </LStack>
  );
}

function FeedBlockRender({ block }: { block: FeedBlock }) {
  let content: ReactNode;

  switch (block.type) {
    case "title":
      content = <FeedTitleBlock />;
      break;
    case "subtitle":
      content = <FeedSubtitleBlock />;
      break;
    case "content":
      content = <FeedContentBlock />;
      break;
    case "cover":
      content = <FeedCoverBlock />;
      break;
    case "categories":
      content = <FeedCategoriesBlock block={block} />;
      break;
    case "threads":
      content = <FeedThreadsBlock block={block} />;
      break;
    case "quick-share":
      content = <FeedQuickShareBlock block={block} />;
      break;
    case "library":
      content = <FeedLibraryBlock block={block} />;
      break;
  }

  return (
    <Box
      className={"feed-page__block feed-page__block--" + block.type}
      data-sd-block={block.type}
      width="full"
    >
      {content}
    </Box>
  );
}

function FeedBlockEditable({
  block,
  index,
}: {
  block: FeedBlock;
  index: number;
}) {
  const {
    attributes,
    listeners,
    setActivatorNodeRef,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: getFeedBlockSortableID(block.type),
    data: {
      type: "feed-block",
      block: block.type,
    } satisfies DragItemFeedBlock,
  });
  const [isOpen, setOpen] = useState(false);

  return (
    <BlockEditor.Root
      ref={setNodeRef}
      className="group"
      data-block-type={block.type}
      style={{
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
      }}
    >
      <BlockEditor.Gutter>
        <BlockEditor.MenuHandle
          {...attributes}
          {...listeners}
          ref={setActivatorNodeRef}
          dragging={isDragging}
          open={isOpen}
          onOpenChange={setOpen}
        >
          <BlockMenu
            block={block}
            index={index}
            onConfigured={() => setOpen(false)}
          />
        </BlockEditor.MenuHandle>
      </BlockEditor.Gutter>
      <BlockEditor.Content>
        <FeedBlockRender block={block} />
      </BlockEditor.Content>
    </BlockEditor.Root>
  );
}

function getFeedBlockSortableID(type: FeedBlockType) {
  return `feed-block:${type}`;
}
