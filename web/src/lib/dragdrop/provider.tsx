"use client";

import {
  type CollisionDetection,
  DndContext,
  type DragEndEvent,
  DragOverlay,
  type DragStartEvent,
  type Modifier,
  PointerSensor,
  TouchSensor,
  pointerWithin,
  rectIntersection,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { SortableData } from "@dnd-kit/sortable";
import { getEventCoordinates } from "@dnd-kit/utilities";
import { useState } from "react";

import { handle } from "@/api/client";
import { Category, Identifier, NodeWithChildren } from "@/api/openapi-schema";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { cx } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { treeView } from "@/styled-system/recipes";

import { useEmitCategoryEvent } from "../category/events";
import { useEmitFeedBlockEvent } from "../feed/events";
import { useEmitLibraryBlockEvent } from "../library/events";
import { useLibraryMutation } from "../library/library";
import { LibraryPageBlockType } from "../library/metadata";
import { useEmitNavigationItemEvent } from "../navigation/events";

export type DragItemLibraryBlock = {
  type: "library-block";
  node: NodeWithChildren; // TODO: Change this to only rely on the node ID.
  block: LibraryPageBlockType;
};

export type DragItemFeedBlock = {
  type: "feed-block";
  block: import("@/lib/settings/feed").FeedBlockType;
};

export type DragItemNavigationItem = {
  type: "navigation-item";
  itemKey: string;
  label: string;
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
  | DragItemLibraryBlock
  | DragItemFeedBlock
  | DragItemNavigationItem
  | DragItemCategory
  | DragItemCategoryDivider;

// The preview is smaller than its source, so preserve sensor movement while
// centering the measured preview on the pointer instead of the source box.
const centerDragOverlayOnPointer: Modifier = ({
  activatorEvent,
  activeNodeRect,
  overlayNodeRect,
  transform,
}) => {
  if (!activatorEvent || !activeNodeRect) {
    return transform;
  }

  const activatorCoordinates = getEventCoordinates(activatorEvent);
  if (!activatorCoordinates) {
    return transform;
  }

  return {
    ...transform,
    x:
      activatorCoordinates.x +
      transform.x -
      activeNodeRect.left -
      (overlayNodeRect?.width ?? activeNodeRect.width) / 2,
    y:
      activatorCoordinates.y +
      transform.y -
      activeNodeRect.top -
      (overlayNodeRect?.height ?? activeNodeRect.height) / 2,
  };
};

const navigationDragOverlayModifiers = [centerDragOverlayOnPointer];

// Ark menu triggers intentionally expose pointer events but not mouse events.
// Keep the dedicated delayed touch sensor while activating mouse drags through
// the pointer event that reaches both plain buttons and composed menu triggers.
class MousePointerSensor extends PointerSensor {
  static override activators: typeof PointerSensor.activators = [
    {
      eventName: "onPointerDown",
      handler: ({ nativeEvent: event }, { onActivation }) => {
        if (
          !event.isPrimary ||
          event.button !== 0 ||
          event.pointerType !== "mouse"
        ) {
          return false;
        }

        onActivation?.({ event });
        return true;
      },
    },
  ];
}

const collisionDetectionStrategy: CollisionDetection = (args) => {
  const activeData = args.active.data.current as DragItemData | undefined;

  if (
    activeData?.type === "feed-block" ||
    activeData?.type === "library-block" ||
    activeData?.type === "navigation-item"
  ) {
    const matchingContainers = args.droppableContainers.filter(
      (container) =>
        (container.data.current as DragItemData | undefined)?.type ===
        activeData.type,
    );
    const filteredArgs = {
      ...args,
      droppableContainers: matchingContainers,
    };
    const pointerCollisions = pointerWithin(filteredArgs);
    return pointerCollisions.length > 0
      ? pointerCollisions
      : rectIntersection(filteredArgs);
  }

  return pointerWithin(args);
};

export function DndProvider({ children }: { children: React.ReactNode }) {
  const { moveNode, revalidate } = useLibraryMutation();
  const emitLibraryBlockEvent = useEmitLibraryBlockEvent();
  const emitFeedBlockEvent = useEmitFeedBlockEvent();
  const emitNavigationItemEvent = useEmitNavigationItemEvent();
  const emitCategoryEvent = useEmitCategoryEvent();
  const [activeItem, setActiveItem] = useState<DragItemData | null>(null);

  const sensors = useSensors(
    useSensor(MousePointerSensor, {
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

  const onDragStart = (event: DragStartEvent) => {
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
    setActiveItem(null);

    if (event.over == null) {
      return;
    }

    // A lot of things are draggable, so each draggable and droppable satisfies
    // a discriminated union of DragItemData.
    const activeData = event.active.data.current as DragItemData;
    const targetData = event.over.data.current as DragItemData;

    if (activeData.type === "node") {
      const active = activeData as DragItemNode & SortableData;
      const target = targetData as DragItemData & SortableData;
      if (target.type !== "node" && target.type !== "divider") {
        return;
      }

      await handle(
        async () => {
          const result = (() => {
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
                    target.type === "divider"
                      ? target.parentID
                      : target.node.id,
                };

              case "node-children": {
                if (active.sortable === undefined) {
                  throw new Error(
                    "Active item is not within a sortable context. This is a bug in Storyden.",
                  );
                }

                if (target.sortable === undefined) {
                  throw new Error(
                    "Target item is not within a sortable context. This is a bug in Storyden.",
                  );
                }

                // NOTE: Is always sortable, for now. May not be in future.
                // When dragging within the same parent, we need to determine if
                // we're moving the item to a position before or after the drop
                // target. If the active item's index is greater than the
                // target's index, we're moving backwards, so insert "above"
                // (before) the target. Otherwise, we're moving forwards, so
                // insert "below" (after) the target.
                const activeIndex = active.sortable.index;
                const targetIndex = target.sortable.index;

                // Safety check: if we don't have valid indices, throw a user-friendly error
                if (activeIndex === -1 || targetIndex === -1) {
                  throw new Error(
                    "Unable to move page: invalid drag position. This is a bug in Storyden.",
                  );
                }

                const movingBackwards = activeIndex > targetIndex;

                return {
                  direction: movingBackwards
                    ? ("above" as const)
                    : ("below" as const),
                  relativeToNode: target.node.id,
                  newParentNode: undefined, // For directory drags, keep in same parent.
                };
              }
            }
          })();

          const { direction, relativeToNode, newParentNode } = result;

          const oldParentID = target.parentID ?? undefined;

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

    if (activeData.type === "library-block") {
      const active = activeData as DragItemLibraryBlock;

      if (targetData.type !== "library-block") {
        return;
      }

      const target = targetData as DragItemLibraryBlock;

      emitLibraryBlockEvent("library:reorder-block", {
        activeId: active.block,
        overId: target.block,
      });
    }

    if (activeData.type === "feed-block") {
      const active = activeData as DragItemFeedBlock;

      if (targetData.type !== "feed-block") {
        return;
      }

      const target = targetData as DragItemFeedBlock;

      emitFeedBlockEvent("feed:reorder-block", {
        activeId: active.block,
        overId: target.block,
      });
    }

    if (activeData.type === "navigation-item") {
      if (targetData.type !== "navigation-item") {
        return;
      }

      emitNavigationItemEvent("navigation:reorder-item", {
        activeKey: activeData.itemKey,
        overKey: targetData.itemKey,
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
      collisionDetection={collisionDetectionStrategy}
      // onDragOver={onDragOver}
      onDragCancel={() => setActiveItem(null)}
      onDragEnd={onDragEnd}
      onDragStart={onDragStart}
    >
      {children}

      <DragOverlay
        adjustScale={false}
        dropAnimation={null}
        className={
          activeItem?.type === "navigation-item"
            ? "navigation-editor__drag-preview-container"
            : undefined
        }
        modifiers={
          activeItem?.type === "navigation-item"
            ? navigationDragOverlayModifiers
            : undefined
        }
      >
        {activeItem &&
          (activeItem.type === "node" ||
            activeItem.type === "category" ||
            activeItem.type === "navigation-item") && (
            <DragOverlaySwitch activeItem={activeItem} />
          )}
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

    case "navigation-item":
      return <DragOverlayNavigationItem activeItem={activeItem} />;

    default:
      return null;
  }
}

function DragOverlayNavigationItem({
  activeItem,
}: {
  activeItem: DragItemNavigationItem;
}) {
  return (
    <div aria-hidden="true" className="navigation-editor__drag-preview">
      {activeItem.label}
    </div>
  );
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
