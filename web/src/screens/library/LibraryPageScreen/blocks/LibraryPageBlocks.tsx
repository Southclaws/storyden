import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { DragItemNodeBlock } from "@/lib/dragdrop/provider";
import { useLibraryBlockEvent } from "@/lib/library/events";
import {
  LibraryPageBlock,
  LibraryPageBlockName,
  LibraryPageBlockType,
} from "@/lib/library/metadata";
import { Box, HStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";

import { useLibraryPageContext } from "../Context";
import { useWatch } from "../store";
import { useEditState } from "../useEditState";

import { BlockMenu } from "./BlockMenu";
import { LibraryPageAssetsBlock } from "./LibraryPageAssetsBlock/LibraryPageAssetsBlock";
import { LibraryPageContentBlock } from "./LibraryPageContentBlock/LibraryPageContentBlock";
import { LibraryPageCoverBlock } from "./LibraryPageCoverBlock/LibraryPageCoverBlock";
import { LibraryPageLinkBlock } from "./LibraryPageLinkBlock/LibraryPageLinkBlock";
import { LibraryPagePropertiesBlock } from "./LibraryPagePropertiesBlock/LibraryPagePropertiesBlock";
import { LibraryPageTableBlock } from "./LibraryPageTableBlock/LibraryPageTableBlock";
import { LibraryPageTagsBlock } from "./LibraryPageTagsBlock/LibraryPageTagsBlock";
import { LibraryPageTitleBlock } from "./LibraryPageTitleBlock/LibraryPageTitleBlock";

export function LibraryPageBlocks() {
  const { store } = useLibraryPageContext();
  const { moveBlock, addBlock, removeBlock } = store.getState();
  const { editing } = useEditState();

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
    (type: LibraryPageBlockType) => {
      addBlock(type);
    },
    [addBlock],
  );
  useLibraryBlockEvent("library:add-block", ({ type }) => {
    handleAddBlock(type);
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

  const blockIds = blocks.map((block) => block.type);

  if (editing) {
    const editStateBlocks = meta.layout?.blocks ?? [];

    return (
      <SortableContext items={blockIds} strategy={verticalListSortingStrategy}>
        {editStateBlocks.map((block) => {
          return <LibraryPageBlockEditable key={block.type} block={block} />;
        })}
      </SortableContext>
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
    case "table":
      return <LibraryPageTableBlock />;
    case "content":
      return <LibraryPageContentBlock />;
  }
}

function LibraryPageBlockEditable({ block }: { block: LibraryPageBlock }) {
  const { initialNode } = useLibraryPageContext();
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: block.type,
    data: {
      type: "block",
      node: initialNode, // TODO: Change this to only pass the node ID.
      block: block.type,
    } as DragItemNodeBlock,
  });

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    maxWidth: "var(--width-adjusted)",
    flexShrink: 0,
    "--width-adjusted": "calc(100% + var(--spacing-5))",
  };

  const dragHandleStyle = {
    cursor: isDragging ? "grabbing" : "grab",
  };

  return (
    <VStack
      className="group"
      style={dragStyle}
      ref={setNodeRef}
      w="full"
      gap="1"
      outlineWidth="thin"
      outlineColor="accent.300"
      outlineStyle="dashed"
      outlineOffset="1"
      borderRadius="sm"
    >
      <WStack>
        <HStack {...attributes} {...listeners} style={dragHandleStyle} gap="1">
          <DragHandleIcon h="4" w="4" color="fg.subtle" />
          <styled.p fontWeight="medium" color="fg.subtle">
            {LibraryPageBlockName[block.type]}
          </styled.p>
        </HStack>
        <BlockMenu block={block} />
      </WStack>
      <Box w="full" minW="0">
        <LibraryPageBlockRender block={block} />
      </Box>
    </VStack>
  );
}
