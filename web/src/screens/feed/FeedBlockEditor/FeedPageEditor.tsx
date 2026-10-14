import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ReactNode, useCallback, useState } from "react";

import { BannerEditor } from "@/components/admin/BrandSettings/BannerEditor/BannerEditor";
import * as BlockEditor from "@/components/ui/block-editor";
import { Button } from "@/components/ui/button";
import { HeadingInput } from "@/components/ui/heading-input";
import { AddIcon } from "@/components/ui/icons/Add";
import { Input } from "@/components/ui/input";
import { DragItemFeedBlock } from "@/lib/dragdrop/provider";
import { useFeedBlockEvent } from "@/lib/feed/events";
import { FeedBlock, FeedBlockType } from "@/lib/settings/feed";
import { LStack } from "@/styled-system/jsx";

import { BlockMenu } from "./BlockMenu";
import { useFeedBlockEditor } from "./Context";
import { CreateBlockMenu } from "./CreateBlockMenu";
import { FeedBlockRender } from "./FeedBlockRender";
import { useDebouncedSiteUpdate } from "./useDebouncedSiteUpdate";

export function FeedPageEditor() {
  const { feed, moveBlock } = useFeedBlockEditor();

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
        <FeedBlockEditableContent block={block} />
      </BlockEditor.Content>
    </BlockEditor.Root>
  );
}

function getFeedBlockSortableID(type: FeedBlockType) {
  return `feed-block:${type}`;
}

function FeedBlockEditableContent({ block }: { block: FeedBlock }) {
  let content: ReactNode;
  switch (block.type) {
    case "title":
      content = <FeedTitleBlockEditor />;
      break;
    case "subtitle":
      content = <FeedSubtitleBlockEditor />;
      break;
    case "cover":
      content = <BannerEditor />;
      break;
    default:
      return <FeedBlockRender block={block} />;
  }
  return <FeedBlockRender block={block}>{content}</FeedBlockRender>;
}

function FeedTitleBlockEditor() {
  const { initialSettings } = useFeedBlockEditor();
  const update = useDebouncedSiteUpdate();
  const title = initialSettings?.title ?? "Storyden";
  return (
    <HeadingInput
      aria-label="Site title"
      defaultValue={title}
      placeholder="Site title"
      onValueChange={(value) => update({ title: value })}
    />
  );
}

function FeedSubtitleBlockEditor() {
  const { initialSettings } = useFeedBlockEditor();
  const update = useDebouncedSiteUpdate();
  const description = initialSettings?.description ?? "";
  return (
    <Input
      aria-label="Site subtitle"
      defaultValue={description}
      placeholder="Site subtitle"
      variant="ghost"
      onChange={(event) => update({ description: event.currentTarget.value })}
    />
  );
}
