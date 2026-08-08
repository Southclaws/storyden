import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback, useState } from "react";

import * as BlockEditor from "@/components/ui/block-editor";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
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
      <LStack gap="4">
        {feed.blocks.map((block) => (
          <FeedBlockRender key={block.type} block={block} />
        ))}
      </LStack>
    );
  }

  return (
    <LStack gap="4">
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
  switch (block.type) {
    case "title":
      return <FeedTitleBlock />;
    case "subtitle":
      return <FeedSubtitleBlock />;
    case "content":
      return <FeedContentBlock />;
    case "cover":
      return <FeedCoverBlock />;
    case "categories":
      return <FeedCategoriesBlock block={block} />;
    case "threads":
      return <FeedThreadsBlock block={block} />;
    case "quick-share":
      return <FeedQuickShareBlock block={block} />;
    case "library":
      return <FeedLibraryBlock block={block} />;
  }
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
        <BlockEditor.Handle>
          <IconButton
            {...attributes}
            {...listeners}
            ref={setActivatorNodeRef}
            aria-label="Move or configure block"
            style={{ cursor: isDragging ? "grabbing" : "grab" }}
            variant="subtle"
            onClick={() => setOpen((open) => !open)}
          >
            <DragHandleIcon />
          </IconButton>

          <Box position="absolute" inset="0" pointerEvents="none">
            <BlockMenu
              block={block}
              index={index}
              open={isOpen}
              onOpenChange={setOpen}
            >
              <Box width="full" height="full" />
            </BlockMenu>
          </Box>
        </BlockEditor.Handle>
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
