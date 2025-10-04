"use client";

import { TreeView as ArkTreeView } from "@ark-ui/react/tree-view";
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

import { handle } from "@/api/client";
import { Category, Identifier, NodeWithChildren } from "@/api/openapi-schema";
import { CategoryBadge } from "@/components/category/CategoryBadge";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { css, cx } from "@/styled-system/css";
import { Box, CardBox, styled } from "@/styled-system/jsx";
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
      const target = targetData as DragItemData;
      if (target.type !== "node" && target.type !== "divider") {
        return;
      }

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
