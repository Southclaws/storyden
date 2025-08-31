import { Portal } from "@ark-ui/react";
import { useDndContext } from "@dnd-kit/core";
import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback, useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import * as Tooltip from "@/components/ui/tooltip";
import { DragItemNodeBlock } from "@/lib/dragdrop/provider";
import { useLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockType } from "@/lib/library/metadata";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";

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

  const blockIds = blocks.map((block) => block.type);

  if (editing) {
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
            <Button variant="outline" size="xs" w="full">
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

  // Manage the menu state manually due to the complexity of the menu trigger
  // also being a drag handle for the block.
  const [isOpen, setOpen] = useState(false);
  function handleMenuToggle() {
    setOpen((prev) => !prev);
  }

  // Manually handle click-away behaviour - the default menu behaviour has been
  // overridden by making it a controlled component in order to allow for the
  // drag handle to be used as a menu open trigger.
  const elementRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!isOpen) return;

    function handleClickAway(event: MouseEvent) {
      if (
        elementRef.current &&
        !elementRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    }

    document.addEventListener("click", handleClickAway);
    return () => document.removeEventListener("click", handleClickAway);
  }, [isOpen]);

  // Check if we're dragging anything at all, to hide the tooltip.
  const { active } = useDndContext();
  const isDraggingAnything = active !== null;

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    flexShrink: 0,
  };

  const dragHandleStyle = {
    cursor: isDragging ? "grabbing" : "grab",
  };

  return (
    <HStack
      id={`block-${block.type}_container`}
      className="group"
      style={dragStyle}
      w="full"
      gap="0"
      position="relative"
    >
      <VStack
        id={`block-${block.type}_gutter-container`}
        ref={setNodeRef}
        w="6"
        left={{ base: "0", md: "-7" }}
        top={{ base: "-7", md: "0" }}
        alignItems="start"
        height="full"
        position="absolute"
        p="0"
      >
        <VStack
          id={`block-${block.type}_gutter-drag-handle`}
          {...listeners}
          {...attributes}
          ref={elementRef}
          w="full"
          h={{ base: "6", md: "full" }}
          bgColor="bg.muted/50"
          color="fg.subtle"
          borderRadius="sm"
          visibility="hidden"
          _groupHover={{
            visibility: "visible",
          }}
          // Hide on mobile: Not happy with the mobile experience of this yet.
          display={{ base: "none", md: "flex" }}
        >
          <Tooltip.Root
            openDelay={0}
            closeDelay={0}
            disabled={isDraggingAnything}
            positioning={{
              slide: true,
              gutter: 4,
              placement: "bottom-start",
            }}
          >
            <Tooltip.Trigger asChild>
              <Box position="relative">
                <Box style={dragHandleStyle}>
                  <IconButton
                    style={dragHandleStyle}
                    id={`block-${block.type}_gutter-drag-handle-button`}
                    variant={{
                      base: "subtle",
                      md: "ghost",
                    }}
                    size="xs"
                    minWidth="5"
                    width="5"
                    height="5"
                    padding="0"
                    color="fg.muted"
                    onClick={handleMenuToggle}
                  >
                    <DragHandleIcon width="4" />
                  </IconButton>
                </Box>
              </Box>
            </Tooltip.Trigger>
            <Portal>
              <Tooltip.Positioner>
                <Tooltip.Arrow>
                  <Tooltip.ArrowTip />
                </Tooltip.Arrow>

                <Tooltip.Content p="1" borderRadius="sm">
                  <p>
                    <styled.span fontWeight="semibold">Click</styled.span>&nbsp;
                    <styled.span fontWeight="normal">to open menu</styled.span>
                  </p>
                  <p>
                    <styled.span fontWeight="semibold">Drag</styled.span>&nbsp;
                    <styled.span fontWeight="normal">to move</styled.span>
                  </p>
                </Tooltip.Content>
              </Tooltip.Positioner>
            </Portal>
          </Tooltip.Root>

          <Box
            position="absolute"
            top="0"
            width="6"
            height="6"
            pointerEvents="none"
          >
            <BlockMenu block={block} open={isOpen} index={index}>
              <Box position="absolute" width="6" height="6" />
            </BlockMenu>
          </Box>
        </VStack>
      </VStack>
      <Box id={`block-${block.type}_content`} w="full" minW="0">
        <LibraryPageBlockRender block={block} />
      </Box>
    </HStack>
  );
}
