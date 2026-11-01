import { Portal } from "@ark-ui/react";
import { useSortable } from "@dnd-kit/sortable";
import { useState } from "react";

import { Identifier, NodeWithChildren } from "@/api/openapi-schema";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import * as Tooltip from "@/components/ui/tooltip";
import { Box, styled } from "@/styled-system/jsx";

type Props = Pick<
  ReturnType<typeof useSortable>,
  "attributes" | "listeners" | "setActivatorNodeRef"
> & {
  node: NodeWithChildren;
  parentID: Identifier;
  overlay: boolean;
  isDragging: boolean;
  isDraggingAnything: boolean;
};

export function DirectoryCardControls({
  node,
  parentID,
  overlay,
  isDragging,
  isDraggingAnything,
  attributes,
  listeners,
  setActivatorNodeRef,
}: Props) {
  const [isOpen, setOpen] = useState(false);
  function handleMenuToggle(event: React.MouseEvent) {
    event.stopPropagation();
    setOpen((open) => !open);
  }

  let cursor: "grab" | "grabbing" = "grab";

  if (isDragging) {
    cursor = "grabbing";
  }

  const dragHandleStyle = { cursor };

  return (
    <Box
      position="absolute"
      top="2"
      right="2"
      zIndex="[6]"
      data-overlay={overlay}
      css={{ "&[data-overlay=true]": { top: "3", right: "3" } }}
    >
      <Tooltip.Root
        openDelay={0}
        closeDelay={0}
        disabled={isDraggingAnything}
        positioning={{
          slide: true,
          gutter: 4,
          placement: "top",
        }}
      >
        <Tooltip.Trigger asChild>
          <Box position="relative" display="flex">
            <IconButton
              ref={setActivatorNodeRef}
              {...listeners}
              {...attributes}
              style={dragHandleStyle}
              variant="subtle"
              minWidth="5"
              width="5"
              height="5"
              data-overlay={overlay}
              css={{
                "&[data-overlay=true]": {
                  minWidth: "6",
                  width: "6",
                  height: "6",
                },
              }}
              padding="0"
              onClick={handleMenuToggle}
              aria-label={`Move ${node.name}`}
              title={`Drag to move; click for page options for ${node.name}`}
            >
              <DragHandleIcon width="4" />
            </IconButton>
          </Box>
        </Tooltip.Trigger>
        <Portal>
          <Tooltip.Positioner>
            <Tooltip.Arrow>
              <Tooltip.ArrowTip />
            </Tooltip.Arrow>

            <Tooltip.Content p="1" borderRadius="sm">
              <p>
                <styled.span fontWeight="semibold">Click</styled.span>
                &nbsp;
                <styled.span fontWeight="normal">to open menu</styled.span>
              </p>
              <p>
                <styled.span fontWeight="semibold">Drag</styled.span>
                &nbsp;
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
        <LibraryPageMenu
          variant="ghost"
          node={node}
          parentID={parentID}
          open={isOpen}
          onOpenChange={(details) => setOpen(details.open)}
          onInteractOutside={() => setOpen(false)}
        >
          <Box position="absolute" width="6" height="6" />
        </LibraryPageMenu>
      </Box>
    </Box>
  );
}
