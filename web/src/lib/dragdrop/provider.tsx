"use client";

import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  MouseSensor,
  TouchSensor,
  pointerWithin,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { SortableData } from "@dnd-kit/sortable";
import { useState } from "react";

import { handle } from "@/api/client";
import { Category, Identifier, NodeWithChildren } from "@/api/openapi-schema";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { cx } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { treeView } from "@/styled-system/recipes";

import { useEmitCategoryEvent } from "../category/events";
import { useEmitLibraryBlockEvent } from "../library/events";
import { useLibraryMutation } from "../library/library";
import { LibraryPageBlockType } from "../library/metadata";

export type DragItemNodeBlock = {
  type: "block";
  node: NodeWithChildren; // TODO: Change this to only rely on the node ID.
  block: LibraryPageBlockType;
};

export type DragItemNode = {
  type: "node";
  node: NodeWithChildren;
  parentID: Identifier | null;
  context: "sidebar" | "node-children";
};

export type DragItemCategory = {
  type: "category";
  categoryID: Identifier;
  category: Category;
  hasChildren: boolean;
};

export type DragItemCategoryDivider = {
  type: "category-divider";
  parentID: Identifier | null;
  siblingCategoryID: Identifier;
  direction: "above" | "below";
};

export type DragItemDivider = {
  type: "divider";
  parentID: Identifier | null;
  siblingNode: NodeWithChildren;
  direction: "above" | "below";
  context: "sidebar";
};

export type DragItemData =
  | DragItemNode
  | DragItemDivider
  | DragItemNodeBlock
  | DragItemCategory
  | DragItemCategoryDivider;

export function DndProvider({ children }: { children: React.ReactNode }) {
  const { moveNode, revalidate } = useLibraryMutation();
  const emitLibraryBlockEvent = useEmitLibraryBlockEvent();
  const emitCategoryEvent = useEmitCategoryEvent();
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
      const target = targetData as DragItemData & SortableData;
      if (target.type !== "node" && target.type !== "divider") {
        return;
      }

      const { direction, relativeToNode, newParentNode } = (() => {
        switch (target.context) {
          case "sidebar":
            return {
              direction:
                target.type === "divider"
                  ? target.direction
                  : ("inside" as const),
              relativeToNode:
                target.type === "divider"
                  ? target.siblingNode.id
                  : target.node.id,
              newParentNode:
                target.type === "divider" ? target.parentID : target.node.id,
            };

          case "node-children": {
            // NOTE: Is always sortable, for now. May not be in future.
            const isTop = target.sortable.index === 0;
            return {
              direction: isTop ? ("above" as const) : ("below" as const),
              relativeToNode: target.node.id,
              newParentNode: undefined, // For directory drags, keep in same parent.
            };
          }
        }
      })();

      const oldParentID = target.parentID ?? undefined;

      await handle(
        async () => {
          await moveNode(
            active.node.id,
            relativeToNode,
            direction,
            newParentNode,
            // NOTE: When a node is dragged between folders in the sidebar,
            // oldParentID is derived from the drop target’s parent instead of
            // the dragged node’s current parent. moveNode uses that ID to
            // optimistically mutate the cached child list for the previous
            // parent. This is intentional as currently, the only use-case for
            // these kinds of drags is moving nodes while looking at the page
            // directory block. This will always be the "old" parent as we want
            // to optimistically revalidate that page (the parent) so its child
            // node API call is revalidated and the new order is rendered.
            // This may need to change in future as this does not currently
            // trigger properly when a member moves a node FROM the sidebar INTO
            // the directory block. In such cases, we will need to revalidate
            // both the old and new parents. Which sounds wasteful, but isn't
            // because there will only ever be a single useSWR call on the page
            // showing a child node list. UNLESS we introduce the ability to
            // show multiple directory blocks on the page with different source
            // parents... but that honestly sounds over-complicated design-wise.
            oldParentID,
          );
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

      if (targetData.type !== "block") {
        return;
      }

      const target = targetData as DragItemNodeBlock;

      emitLibraryBlockEvent("library:reorder-block", {
        activeId: active.block,
        overId: target.block,
      });
    }

    if (activeData.type === "category") {
      const active = activeData as DragItemCategory;
      const target = targetData as DragItemCategoryDivider | DragItemCategory;
      if (
        targetData.type !== "category" &&
        targetData.type !== "category-divider"
      ) {
        return;
      }

      const direction =
        target.type === "category-divider" ? target.direction : "inside";
      const targetCategory =
        target.type === "category-divider"
          ? target.siblingCategoryID
          : target.categoryID;
      const newParent =
        target.type === "category-divider"
          ? target.parentID
          : target.categoryID;

      emitCategoryEvent("category:reorder-category", {
        categorySlug: active.category.slug,
        targetCategory,
        direction,
        newParent,
      });
    }
  };

  return (
    <DndContext
      id="sd-dnd"
      sensors={sensors}
      collisionDetection={pointerWithin}
      // onDragOver={onDragOver}
      onDragEnd={onDragEnd}
      onDragStart={onDragStart}
    >
      {children}

      <DragOverlay>
        {activeItem && <DragOverlaySwitch activeItem={activeItem} />}
      </DragOverlay>
    </DndContext>
  );
}

type DragOverlaySwitchProps = {
  activeItem: DragItemData;
};

function DragOverlaySwitch({ activeItem }: DragOverlaySwitchProps) {
  switch (activeItem.type) {
    case "node":
      return <DragOverlayNavigationNode activeItem={activeItem} />;

    case "category":
      return <DragOverlayNavigationCategory activeItem={activeItem} />;

    default:
      return null;
  }
}

function DragOverlayNavigationCategory({
  activeItem,
}: {
  activeItem: DragItemCategory;
}) {
  const styles = treeView();
  return (
    <Box className={cx(styles.branchControl)} opacity="5">
      <Box className={styles.branchIndicator}>
        {activeItem?.hasChildren ? <ChevronRightIcon /> : <BulletIcon />}
      </Box>

      <Box className={styles.branchText}>{activeItem.category.name}</Box>
    </Box>
  );
}

function DragOverlayNavigationNode({
  activeItem,
}: {
  activeItem: DragItemNode;
}) {
  const styles = treeView();
  return (
    <Box className={cx(styles.branchControl)} opacity="5">
      <Box className={styles.branchIndicator}>
        {activeItem?.node.children.length > 0 ? (
          <ChevronRightIcon />
        ) : (
          <BulletIcon />
        )}
      </Box>

      <Box className={styles.branchText}>{activeItem.node.name}</Box>
    </Box>
  );
}
