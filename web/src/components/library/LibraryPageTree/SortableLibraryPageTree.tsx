"use client";

import {
  DndContext,
  type DragEndEvent,
  DragOverlay,
  type DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ChevronDown, ChevronRight, GripVertical } from "lucide-react";
import * as React from "react";

import { cn } from "@/lib/utils";

// Define the tree item type
interface TreeItemData {
  id: string;
  name: string;
  children?: TreeItemData[];
}

// Initial tree data
const initialItems: TreeItemData[] = [
  {
    id: "1",
    name: "Root 1",
    children: [
      { id: "1-1", name: "Child 1-1" },
      {
        id: "1-2",
        name: "Child 1-2",
        children: [
          { id: "1-2-1", name: "Grandchild 1-2-1" },
          { id: "1-2-2", name: "Grandchild 1-2-2" },
        ],
      },
      { id: "1-3", name: "Child 1-3" },
    ],
  },
  {
    id: "2",
    name: "Root 2",
    children: [
      { id: "2-1", name: "Child 2-1" },
      { id: "2-2", name: "Child 2-2" },
    ],
  },
  {
    id: "3",
    name: "Root 3",
  },
];

// Flatten tree for easier manipulation
const flattenTree = (
  items: TreeItemData[],
  parentId: string | null = null,
): { item: TreeItemData; parentId: string | null }[] => {
  return items.reduce(
    (acc, item) => {
      acc.push({ item, parentId });
      if (item.children) {
        acc.push(...flattenTree(item.children, item.id));
      }
      return acc;
    },
    [] as { item: TreeItemData; parentId: string | null }[],
  );
};

// Custom tree component that doesn't rely on Ark UI's TreeView
const CustomTreeItem = ({
  item,
  depth = 0,
  isExpanded = true,
  onToggle,
  onDragStart,
  expandedItems,
}: {
  item: TreeItemData;
  depth?: number;
  isExpanded?: boolean;
  onToggle: (id: string) => void;
  onDragStart: (id: string) => void;
  expandedItems: string[];
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: item.id,
    data: {
      type: "tree-item",
      item,
      depth,
    },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.4 : 1,
  };

  const hasChildren = item.children && item.children.length > 0;

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn("relative", isDragging ? "z-10" : "")}
    >
      <div
        className="flex items-center py-1"
        style={{ paddingLeft: `${depth * 16 + 8}px` }}
      >
        {hasChildren ? (
          <button
            type="button"
            onClick={() => onToggle(item.id)}
            className="mr-1 flex items-center justify-center w-4 h-4"
          >
            {isExpanded ? (
              <ChevronDown className="h-4 w-4" />
            ) : (
              <ChevronRight className="h-4 w-4" />
            )}
          </button>
        ) : (
          <div className="w-4 h-4 mr-1" />
        )}

        <div className="flex-1 flex items-center">
          <div
            {...attributes}
            {...listeners}
            className="mr-2 cursor-grab active:cursor-grabbing"
            onMouseDown={() => onDragStart(item.id)}
          >
            <GripVertical className="h-4 w-4 text-muted-foreground" />
          </div>
          <span>{item.name}</span>
        </div>
      </div>

      {hasChildren && isExpanded && (
        <div>
          <SortableContext
            items={item.children.map((child) => child.id)}
            strategy={verticalListSortingStrategy}
          >
            {item.children.map((child) => (
              <CustomTreeItem
                key={child.id}
                item={child}
                depth={depth + 1}
                isExpanded={expandedItems.includes(child.id)}
                onToggle={onToggle}
                onDragStart={onDragStart}
                expandedItems={expandedItems}
              />
            ))}
          </SortableContext>
        </div>
      )}
    </div>
  );
};

// Main draggable tree component
export function DraggableTree() {
  const [items, setItems] = React.useState<TreeItemData[]>(initialItems);
  const [activeId, setActiveId] = React.useState<string | null>(null);
  const [expandedItems, setExpandedItems] = React.useState<string[]>([
    "1",
    "1-2",
  ]);

  // Find the active item for the drag overlay
  const activeItem = React.useMemo(() => {
    if (!activeId) return null;
    const flatItems = flattenTree(items);
    return flatItems.find(({ item }) => item.id === activeId)?.item || null;
  }, [activeId, items]);

  // Configure DnD sensors
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    }),
  );

  // Handle drag start
  const handleDragStart = (event: DragStartEvent) => {
    const { id } = event.active;
    setActiveId(id.toString());
  };

  // Handle drag end
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      setItems((prevItems) => {
        // Find the items we're working with
        const flatItems = flattenTree(prevItems);
        const activeItemWithParent = flatItems.find(
          ({ item }) => item.id === active.id,
        );
        const overItemWithParent = flatItems.find(
          ({ item }) => item.id === over.id,
        );

        if (!activeItemWithParent || !overItemWithParent) return prevItems;

        const activeParentId = activeItemWithParent.parentId;
        const overParentId = overItemWithParent.parentId;

        // Create a new tree with the item removed
        const removeItem = (
          items: TreeItemData[],
          id: string,
        ): [TreeItemData[], TreeItemData | null] => {
          let removed: TreeItemData | null = null;

          const newItems = items.reduce((acc, item) => {
            if (item.id === id) {
              removed = { ...item };
              return acc;
            }

            if (item.children) {
              const [newChildren, removedFromChildren] = removeItem(
                item.children,
                id,
              );
              if (removedFromChildren) {
                removed = removedFromChildren;
                return [...acc, { ...item, children: newChildren }];
              }
            }

            return [...acc, item];
          }, [] as TreeItemData[]);

          return [newItems, removed];
        };

        // Insert an item at a specific position
        const insertItem = (
          items: TreeItemData[],
          parentId: string | null,
          item: TreeItemData,
          overId: string,
        ): TreeItemData[] => {
          if (parentId === null) {
            // Insert at root level
            const index = items.findIndex((i) => i.id === overId);
            if (index === -1) return [...items, item];

            const newItems = [...items];
            newItems.splice(index, 0, item);
            return newItems;
          }

          // Insert at nested level
          return items.map((currentItem) => {
            if (currentItem.id === parentId) {
              const children = currentItem.children || [];
              const index = children.findIndex((child) => child.id === overId);

              if (index === -1) {
                return {
                  ...currentItem,
                  children: [...children, item],
                };
              }

              const newChildren = [...children];
              newChildren.splice(index, 0, item);

              return {
                ...currentItem,
                children: newChildren,
              };
            }

            if (currentItem.children) {
              return {
                ...currentItem,
                children: insertItem(
                  currentItem.children,
                  parentId,
                  item,
                  overId,
                ),
              };
            }

            return currentItem;
          });
        };

        // Remove the active item from its current position
        const [itemsWithoutActive, removedItem] = removeItem(
          prevItems,
          active.id.toString(),
        );
        if (!removedItem) return prevItems;

        // Insert the active item at its new position
        return insertItem(
          itemsWithoutActive,
          overParentId,
          removedItem,
          over.id.toString(),
        );
      });
    }

    setActiveId(null);
  };

  // Toggle expanded state
  const toggleExpanded = (id: string) => {
    setExpandedItems((prev) =>
      prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id],
    );
  };

  // Get all the IDs for the sortable context
  const rootIds = items.map((item) => item.id);

  return (
    <div className="w-full max-w-md mx-auto p-4 border rounded-lg bg-card">
      <h2 className="text-xl font-semibold mb-4">Draggable Tree</h2>

      <DndContext
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className="w-full">
          <SortableContext
            items={rootIds}
            strategy={verticalListSortingStrategy}
          >
            {items.map((item) => (
              <CustomTreeItem
                key={item.id}
                item={item}
                isExpanded={expandedItems.includes(item.id)}
                onToggle={toggleExpanded}
                onDragStart={setActiveId}
                expandedItems={expandedItems}
              />
            ))}
          </SortableContext>
        </div>

        <DragOverlay>
          {activeId && activeItem && (
            <div className="bg-background border rounded-md shadow-md p-2 w-64">
              <div className="flex items-center">
                <GripVertical className="h-4 w-4 mr-2 text-muted-foreground" />
                <span>{activeItem.name}</span>
              </div>
            </div>
          )}
        </DragOverlay>
      </DndContext>
    </div>
  );
}
