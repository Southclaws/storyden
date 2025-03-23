import { TreeView as ArkTreeView } from "@ark-ui/react/tree-view";
import {
  DndContext,
  DragEndEvent,
  DragMoveEvent,
  DragOverlay,
  DragStartEvent,
  MouseSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useCallback, useState } from "react";

import { NodeWithChildren } from "@/api/openapi-schema";

import {
  LibraryPageTree,
  TreeViewData,
  TreeViewProps,
} from "./LibraryPageTree";

interface DraggableLibraryPageTreeProps extends TreeViewProps {
  onReorder?: (sourceId: string, targetId: string) => void;
  onMove?: (sourceId: string, targetId: string) => void;
}

export function DraggableLibraryPageTree({
  data,
  onReorder,
  onMove,
  ...props
}: DraggableLibraryPageTreeProps) {
  const [activeId, setActiveId] = useState<string | null>(null);
  const [overId, setOverId] = useState<string | null>(null);

  const sensors = useSensors(
    useSensor(MouseSensor, {
      activationConstraint: {
        distance: 8, // Start dragging after moving 8px
      },
    }),
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
  );

  const findNode = useCallback(
    (
      id: string,
      nodes: NodeWithChildren[] = data.children,
    ): NodeWithChildren | null => {
      for (const node of nodes) {
        if (node.id === id) return node;
        if (node.children) {
          const found = findNode(id, node.children);
          if (found) return found;
        }
      }
      return null;
    },
    [data],
  );

  const findParentNode = useCallback(
    (
      id: string,
      nodes: NodeWithChildren[] = data.children,
      parent: NodeWithChildren | null = null,
    ): { parent: NodeWithChildren | null; index: number } => {
      for (let i = 0; i < nodes.length; i++) {
        if (nodes[i].id === id) {
          return { parent, index: i };
        }
        if (nodes[i].children) {
          const found = findParentNode(id, nodes[i].children, nodes[i]);
          if (found.parent !== null) {
            return found;
          }
        }
      }
      return { parent: null, index: -1 };
    },
    [data],
  );

  const handleDragStart = ({ active }: DragStartEvent) => {
    setActiveId(active.id as string);
  };

  const handleDragMove = ({ over }: DragMoveEvent) => {
    setOverId((over?.id as string) ?? null);
  };

  const handleDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over) return;

    const activeId = active.id as string;
    const overId = over.id as string;

    if (activeId === overId) return;

    const activeNode = findNode(activeId);
    const overNode = findNode(overId);

    if (!activeNode || !overNode) return;

    const { parent: activeParent } = findParentNode(activeId);
    const { parent: overParent } = findParentNode(overId);

    if (activeParent === overParent) {
      // Reordering within the same parent
      if (onReorder) {
        onReorder(activeId, overId);
      }
    } else {
      // Moving to a different parent
      if (onMove) {
        onMove(activeId, overId);
      }
    }

    setActiveId(null);
    setOverId(null);
  };

  const handleDragCancel = () => {
    setActiveId(null);
    setOverId(null);
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={handleDragStart}
      onDragMove={handleDragMove}
      onDragEnd={handleDragEnd}
      onDragCancel={handleDragCancel}
    >
      <SortableContext
        items={data.children?.map((node) => node.id)}
        strategy={verticalListSortingStrategy}
      >
        <LibraryPageTree data={data} {...props} />
      </SortableContext>

      <DragOverlay>
        {activeId ? (
          <div style={{ opacity: 0.8 }}>{findNode(activeId)?.name}</div>
        ) : null}
      </DragOverlay>
    </DndContext>
  );
}
