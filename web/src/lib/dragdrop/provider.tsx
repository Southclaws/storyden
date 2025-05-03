"use client";

import {
  DndContext,
  DragEndEvent,
  DragOverEvent,
  MouseSensor,
  TouchSensor,
  pointerWithin,
  useSensor,
  useSensors,
} from "@dnd-kit/core";

import { handle } from "@/api/client";
import { Identifier, NodeWithChildren } from "@/api/openapi-schema";

import { useLibraryMutation } from "../library/library";

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

export type DragItemData = DragItemNode | DragItemDivider;

export function DndProvider({ children }: { children: React.ReactNode }) {
  const { moveNode, revalidate } = useLibraryMutation();

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

    const active = event.active.data.current as DragItemNode;
    const target = event.over.data.current as DragItemData;

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
  };

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={pointerWithin}
      // onDragOver={onDragOver}
      onDragEnd={onDragEnd}
    >
      {children}
    </DndContext>
  );
}
