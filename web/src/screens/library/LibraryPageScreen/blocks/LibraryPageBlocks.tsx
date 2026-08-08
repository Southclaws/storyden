import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback, useState } from "react";

import * as BlockEditor from "@/components/ui/block-editor";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DragItemLibraryBlock } from "@/lib/dragdrop/provider";
import { useLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockType } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../Context";
import { useWatch } from "../store";
import { useEditState } from "../useEditState";

import { BlockMenu } from "./BlockMenu";
import { CreateBlockMenu } from "./CreateBlockMenu";
import { LibraryPageAssetsBlock } from "./LibraryPageAssetsBlock/LibraryPageAssetsBlock";
import { LibraryPageContentBlock } from "./LibraryPageContentBlock/LibraryPageContentBlock";
import { LibraryPageCoverBlock } from "./LibraryPageCoverBlock/LibraryPageCoverBlock";
import { LibraryPageDirectoryBlock } from "./LibraryPageDirectoryBlock/LibraryPageDirectoryBlock";
import { LibraryPageLinkBlock } from "./LibraryPageLinkBlock/LibraryPageLinkBlock";
import { LibraryPagePropertiesBlock } from "./LibraryPagePropertiesBlock/LibraryPagePropertiesBlock";
import { LibraryPageTagsBlock } from "./LibraryPageTagsBlock/LibraryPageTagsBlock";
import { LibraryPageTitleBlock } from "./LibraryPageTitleBlock/LibraryPageTitleBlock";

export function LibraryPageBlocks() {
  const { initialNode, store } = useLibraryPageContext();
  const { moveBlock, addBlock, removeBlock } = store.getState();
  const { isDirectEditing } = useEditState();

  const meta = useWatch((s) => s.draft.meta);

  const handleReorder = useCallback(
    (activeId: LibraryPageBlockType, overId: LibraryPageBlockType) => {
      if (!meta.layout) {
        throw new Error("No layout found in metadata");
      }

      const index = meta.layout.blocks.findIndex((b) => b.type === overId);

      moveBlock(activeId, index);
    },
    [moveBlock, meta],
  );
  useLibraryBlockEvent("library:reorder-block", ({ activeId, overId }) => {
    handleReorder(activeId, overId);
  });

  const handleAddBlock = useCallback(
    (type: LibraryPageBlockType, index?: number) => {
      addBlock(type, index);
    },
    [addBlock],
  );
  useLibraryBlockEvent("library:add-block", ({ type, index }) => {
    handleAddBlock(type, index);
  });

  const handleRemoveBlock = useCallback(
    (type: LibraryPageBlockType) => {
      removeBlock(type);
    },
    [removeBlock],
  );
  useLibraryBlockEvent("library:remove-block", ({ type }) => {
    handleRemoveBlock(type);
  });

  const blocks = meta.layout?.blocks ?? [];

  const blockIds = blocks.map((block) =>
    getLibraryBlockSortableID(initialNode.id, block.type),
  );

  if (isDirectEditing) {
    const editStateBlocks = meta.layout?.blocks ?? [];

    return (
      <>
        <SortableContext
          items={blockIds}
          strategy={verticalListSortingStrategy}
        >
          {editStateBlocks.map((block, index) => {
            return (
              <LibraryPageBlockEditable
                key={block.type}
                block={block}
                index={index}
              />
            );
          })}
        </SortableContext>

        <CreateBlockMenu
          trigger={
            <Button variant="outline" w="full">
              <AddIcon />
              &nbsp;Add Block
            </Button>
          }
          positioning={{
            placement: "bottom",
          }}
        />
      </>
    );
  }

  return (
    <>
      {blocks.map((block) => {
        return <LibraryPageBlockRender key={block.type} block={block} />;
      })}
    </>
  );
}

function LibraryPageBlockRender({ block }: { block: LibraryPageBlock }) {
  switch (block.type) {
    case "cover":
      return <LibraryPageCoverBlock />;
    case "assets":
      return <LibraryPageAssetsBlock />;
    case "title":
      return <LibraryPageTitleBlock />;
    case "tags":
      return <LibraryPageTagsBlock />;
    case "link":
      return <LibraryPageLinkBlock />;
    case "properties":
      return <LibraryPagePropertiesBlock />;
    case "directory":
      return <LibraryPageDirectoryBlock />;
    case "content":
      return <LibraryPageContentBlock />;
  }
}

function LibraryPageBlockEditable({
  block,
  index,
}: {
  block: LibraryPageBlock;
  index: number;
}) {
  const { initialNode } = useLibraryPageContext();
  const {
    attributes,
    listeners,
    setActivatorNodeRef,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: getLibraryBlockSortableID(initialNode.id, block.type),
    data: {
      type: "library-block",
      node: initialNode, // TODO: Change this to only pass the node ID.
      block: block.type,
    } as DragItemLibraryBlock,
  });
  const [isOpen, setOpen] = useState(false);

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    flexShrink: 0,
  };

  return (
    <BlockEditor.Root
      ref={setNodeRef}
      id={`block-${block.type}_container`}
      className="group"
      data-block-type={block.type}
      style={dragStyle}
    >
      <BlockEditor.Gutter id={`block-${block.type}_gutter-container`}>
        <BlockEditor.MenuHandle
          {...attributes}
          {...listeners}
          ref={setActivatorNodeRef}
          dragging={isDragging}
          id={`block-${block.type}_gutter-drag-handle`}
          open={isOpen}
          onOpenChange={setOpen}
        >
          <BlockMenu block={block} index={index} />
        </BlockEditor.MenuHandle>
      </BlockEditor.Gutter>
      <BlockEditor.Content id={`block-${block.type}_content`}>
        <LibraryPageBlockRender block={block} />
      </BlockEditor.Content>
    </BlockEditor.Root>
  );
}

function getLibraryBlockSortableID(
  nodeID: string,
  blockType: LibraryPageBlockType,
) {
  return `library-block:${nodeID}:${blockType}`;
}
