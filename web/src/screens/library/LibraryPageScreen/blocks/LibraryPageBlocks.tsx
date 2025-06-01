import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback } from "react";
import { FixedCropperRef } from "react-advanced-cropper";

import { NodeWithChildren } from "@/api/openapi-schema";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { DragItemNodeBlock } from "@/lib/dragdrop/provider";
import { useLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock, LibraryPageBlockType } from "@/lib/library/metadata";
import { Box, HStack, VStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../Context";
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

type Props = {
  cropperRef: React.RefObject<FixedCropperRef | null>;
};

export function LibraryPageBlocks({ cropperRef }: Props) {
  const { node, form } = useLibraryPageContext();
  const { editing } = useEditState();

  const meta = form.watch("meta", node.meta);

  const handleReorder = useCallback(
    (activeId: LibraryPageBlockType, overId: LibraryPageBlockType) => {
      if (!meta.layout) {
        return;
      }

      const currentBlocks = meta.layout.blocks;

      const reOrderBlocks = () => {
        const activeIndex = currentBlocks.findIndex((b) => b.type === activeId);
        const overIndex = currentBlocks.findIndex((b) => b.type === overId);

        if (
          activeIndex === -1 ||
          overIndex === -1 ||
          activeIndex === overIndex
        ) {
          return currentBlocks;
        }

        const newBlocks = [...currentBlocks];
        const [movedBlock] = newBlocks.splice(activeIndex, 1);
        if (!movedBlock) {
          return newBlocks;
        }

        newBlocks.splice(overIndex, 0, movedBlock);

        return newBlocks;
      };

      const newBlocks = reOrderBlocks();

      const newMeta = {
        ...meta,
        layout: {
          ...meta.layout,
          blocks: newBlocks,
        },
      };

      form.setValue("meta", newMeta);
    },
    [form, meta],
  );

  const handleAddBlock = useCallback(
    (type: LibraryPageBlockType) => {
      if (!meta.layout) {
        return;
      }

      const currentBlocks = meta.layout.blocks;

      const newBlocks = [...currentBlocks, { type }] as LibraryPageBlock[];

      const newMeta = {
        ...meta,
        layout: {
          ...meta.layout,
          blocks: newBlocks,
        },
      };

      form.setValue("meta", newMeta);
    },
    [form, meta],
  );

  const handleRemoveBlock = useCallback(
    (type: LibraryPageBlockType) => {
      if (!meta.layout) {
        return;
      }

      const currentBlocks = meta.layout.blocks;

      const newBlocks = currentBlocks.filter(
        (block) => block.type !== type,
      ) as LibraryPageBlock[];

      const newMeta = {
        ...meta,
        layout: {
          ...meta.layout,
          blocks: newBlocks,
        },
      };

      form.setValue("meta", newMeta);
    },
    [form, meta],
  );

  useLibraryBlockEvent("library:reorder-block", ({ activeId, overId }) => {
    handleReorder(activeId, overId);
  });

  useLibraryBlockEvent("library:add-block", ({ type }) => {
    handleAddBlock(type);
  });

  useLibraryBlockEvent("library:remove-block", ({ type }) => {
    handleRemoveBlock(type);
  });

  const blocks = meta.layout?.blocks ?? [];

  const blockIds = blocks.map((block) => block.type);

  if (editing) {
    const editStateBlocks = meta.layout?.blocks ?? [];

    return (
      <SortableContext items={blockIds} strategy={verticalListSortingStrategy}>
        {editStateBlocks.map((block, index) => {
          return (
            <LibraryPageBlockEditable
              key={index}
              cropperRef={cropperRef}
              block={block}
              node={node}
            />
          );
        })}
      </SortableContext>
    );
  }

  return (
    <>
      {blocks.map((block, index) => {
        return (
          <LibraryPageBlockRender
            key={index}
            cropperRef={cropperRef}
            block={block}
          />
        );
      })}
    </>
  );
}

function LibraryPageBlockRender({
  cropperRef,
  block,
}: Props & { block: LibraryPageBlock }) {
  switch (block.type) {
    case "cover":
      return <LibraryPageCoverBlock ref={cropperRef} />;
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

function LibraryPageBlockEditable({
  cropperRef,
  block,
  node,
}: Props & { block: LibraryPageBlock; node: NodeWithChildren }) {
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
      node: node,
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
    <HStack
      className="group"
      style={dragStyle}
      ref={setNodeRef}
      w="var(--width-adjusted)"
      ml="-5"
      gap="0"
    >
      <VStack
        {...attributes}
        {...listeners}
        style={dragHandleStyle}
        w="5"
        pr="1"
        alignItems="start"
        height="full"
        position="relative"
      >
        <VStack
          position="absolute"
          w="full"
          color="fg.subtle"
          borderRadius="sm"
          visibility="hidden"
          _groupHover={{
            bgColor: "bg.muted",
            visibility: "visible",
          }}
          title={block.type}
          gap="1"
        >
          <DragHandleIcon width="4" />
          <BlockMenu node={node} block={block} />
        </VStack>
      </VStack>
      <Box w="full" minW="0">
        <LibraryPageBlockRender cropperRef={cropperRef} block={block} />
      </Box>
    </HStack>
  );
}
