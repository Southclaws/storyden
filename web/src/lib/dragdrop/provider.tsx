"use client";

import {
  DndContext,
  DragEndEvent,
  DragOverEvent,
  DragOverlay,
  MouseSensor,
  TouchSensor,
  pointerWithin,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { useState } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { Identifier, NodeWithChildren } from "@/api/openapi-schema";

import { useEmitLibraryBlockEvent } from "../library/events";
import { useLibraryMutation } from "../library/library";
import { LibraryPageBlockType } from "../library/metadata";

export type DragItemNodeBlock = {
  type: "block";
  node: NodeWithChildren;
  block: LibraryPageBlockType;
};

export type DragItemNode = {
  type: "node";
  node: NodeWithChildren;
};

export type DragItemDivider = {
  type: "divider";
  parentID: Identifier | null;
  siblingNode: NodeWithChildren;
  direction: "above" | "below";
};

export type DragItemData = DragItemNode | DragItemDivider | DragItemNodeBlock;

export function DndProvider({ children }: { children: React.ReactNode }) {
  const { moveNode, revalidate } = useLibraryMutation();
  const emitLibraryBlockEvent = useEmitLibraryBlockEvent();
  const [activeItem, setActiveItem] = useState<DragItemData | null>(null);

  const sensors = useSensors(
    useSensor(MouseSensor, {
      activationConstraint: {
        distance: 2,
      },
    }),
    useSensor(TouchSensor, {
      activationConstraint: {
        delay: 150,
        tolerance: 5,
      },
    }),
  );

  const onDragStart = (event: DragEndEvent) => {
    const activeData = event.active.data.current as DragItemData;
    setActiveItem(activeData);
  };

  // NOTE: Unused currently.
  // const onDragOver = (event: DragOverEvent) => {
  //   const active = event.active.data.current as DragItemNode;
  //   const target = event.over?.data.current as DragItemData;
  //   console.log("onDragOver", active, target);
  // };

  const onDragEnd = async (event: DragEndEvent) => {
    if (event.over == null) {
      return;
    }

    // A lot of things are draggable, so each draggable and droppable satisfies
    // a discriminated union of DragItemData.
    const activeData = event.active.data.current as DragItemData;
    const targetData = event.over.data.current as DragItemData;

    if (activeData.type === "node") {
      const active = activeData as DragItemNode;
      const target = targetData as DragItemData;

      const direction = target.type === "divider" ? target.direction : "inside";
      const targetNode =
        target.type === "divider" ? target.siblingNode : target.node;
      const newParent =
        target.type === "divider" ? target.parentID : target.node.id;

      await handle(
        async () => {
          await moveNode(active.node.id, targetNode.id, direction, newParent);
        },
        {
          async cleanup() {
            await revalidate();
          },
        },
      );
    }

    if (activeData.type === "block") {
      const active = activeData as DragItemNodeBlock;
      const target = targetData as DragItemNodeBlock;

      emitLibraryBlockEvent("library:reorder-block", {
        activeId: active.block,
        overId: target.block,
      });
    }
  };

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={pointerWithin}
      // onDragOver={onDragOver}
      onDragEnd={onDragEnd}
      onDragStart={onDragStart}
    >
      {children}

      <DragOverlay>
        {/* TODO: Drag previews for different element types */}
        {/* <p>{activeItem?.type}</p> */}
      </DragOverlay>
    </DndContext>
  );
}
