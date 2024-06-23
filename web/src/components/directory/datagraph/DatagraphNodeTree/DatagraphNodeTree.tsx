import type { Assign } from "@ark-ui/react";
import {
  TreeView as ArkTreeView,
  type TreeViewRootProps,
} from "@ark-ui/react/tree-view";
import {
  DndContext,
  DragOverlay,
  closestCenter,
  closestCorners,
  useDraggable,
  useDroppable,
} from "@dnd-kit/core";
import { CSS, useCombinedRefs } from "@dnd-kit/utilities";
import Link from "next/link";
import { forwardRef } from "react";

import { NodeWithChildren } from "@/api/openapi/schemas";
import { css, cx } from "@/styled-system/css";
import { splitCssProps } from "@/styled-system/jsx";
import { type TreeViewVariantProps, treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";
import type { JsxStyleProps } from "@/styled-system/types";

import { DatagraphNodeMenu } from "../DatagraphNodeMenu/DatagraphNodeMenu";

import { useDatagraphNodeTree } from "./useDatagraphNodeTree";

export interface TreeViewData {
  label: string;
  children: NodeWithChildren[];
}

export interface TreeViewProps
  extends Assign<JsxStyleProps, TreeViewRootProps>,
    TreeViewVariantProps {
  data: TreeViewData;
}

export const DatagraphNodeTree = forwardRef<HTMLDivElement, TreeViewProps>(
  (props, ref) => {
    const { sensors, handleDragEnd, handleDelete } = useDatagraphNodeTree();
    const [cssProps, localProps] = splitCssProps(props);
    const { data, className, ...rootProps } = localProps;
    const styles = treeView();

    const renderChild = (child: NodeWithChildren) => {
      return (
        <ArkTreeView.Branch
          key={child.id}
          value={child.id}
          className={styles.branch}
        >
          <TreeBranch
            styles={styles}
            child={child}
            handleDelete={handleDelete}
          />

          <ArkTreeView.BranchContent className={styles.branchContent}>
            {child.children?.map((child) =>
              child.children ? (
                renderChild(child)
              ) : (
                <TreeItem
                  key={child.id}
                  styles={styles}
                  child={child}
                  handleDelete={handleDelete}
                />
              ),
            )}
          </ArkTreeView.BranchContent>
        </ArkTreeView.Branch>
      );
    };

    return (
      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragEnd={handleDragEnd}
        onDragStart={console.log}
      >
        <ArkTreeView.Root
          ref={ref}
          aria-label={data.label}
          className={cx(styles.root, css(cssProps), className)}
          {...rootProps}
        >
          <ArkTreeView.Tree className={styles.tree}>
            {data.children.map(renderChild)}
          </ArkTreeView.Tree>
        </ArkTreeView.Root>

        {/* <DragOverlay>
          <p>test</p>
        </DragOverlay> */}
      </DndContext>
    );
  },
);

DatagraphNodeTree.displayName = "DatagraphNodeTree";

type BranchProps = {
  child: NodeWithChildren;
  styles: any;
  handleDelete: (slug: string) => void;
};

function TreeBranch({ styles, child, handleDelete }: BranchProps) {
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
    transform,
    isDragging,
  } = useDraggable({ id: child.id });

  const { setNodeRef: setDroppableNodeRef } = useDroppable({
    id: child.id,
  });

  const setNodeRef = useCombinedRefs(setDraggableNodeRef, setDroppableNodeRef);

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
  };

  const conditionalDragStyles = css({
    cursor: isDragging ? "grabbing" : "grab",
    ...(isDragging && {
      pointerEvents: "none",
    }),
  });

  return (
    <ArkTreeView.BranchControl
      ref={setNodeRef}
      style={dragStyle}
      className={cx(styles.branchControl, conditionalDragStyles)}
      {...attributes}
      {...listeners}
    >
      <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
        {child.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
      </ArkTreeView.BranchIndicator>

      <ArkTreeView.BranchText asChild className={styles.branchText}>
        {isDragging ? (
          <p>{child.name}</p>
        ) : (
          <Link
            href={`/directory/${child.slug}`}
            className={conditionalDragStyles}
          >
            {child.name}
          </Link>
        )}
      </ArkTreeView.BranchText>

      <DatagraphNodeMenu node={child} onDelete={handleDelete} />
    </ArkTreeView.BranchControl>
  );
}

function TreeItem({ styles, child }: BranchProps) {
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
    transform,
    isDragging,
  } = useDraggable({ id: child.id });

  const { setNodeRef: setDroppableNodeRef } = useDroppable({
    id: child.id,
  });

  const setNodeRef = useCombinedRefs(setDraggableNodeRef, setDroppableNodeRef);

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
  };

  return (
    <ArkTreeView.Item
      ref={setNodeRef}
      style={dragStyle}
      value={child.id}
      className={cx(
        styles.item,
        css({
          cursor: isDragging ? "grabbing" : "grab",
          ...(isDragging && {
            pointerEvents: "none",
          }),
        }),
      )}
      {...attributes}
      {...listeners}
    >
      <ArkTreeView.ItemText className={styles.itemText}>
        <Link
          href={`/directory/${child.slug}`}
          className={conditionalDragStyle}
        >
          {child.name}
        </Link>
      </ArkTreeView.ItemText>
    </ArkTreeView.Item>
  );
}

const ChevronRightIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
    <title>Chevron Right Icon</title>
    <path
      fill="none"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="m9 18l6-6l-6-6"
    />
  </svg>
);

const BulletIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill={token("colors.fg.muted")}
  >
    <circle cx="12.1" cy="12.1" r="2.5" />
  </svg>
);
