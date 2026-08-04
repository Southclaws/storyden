"use client";

import {
  TreeView as ArkTreeView,
  type TreeViewExpandedChangeDetails,
  createTreeCollection,
} from "@ark-ui/react/tree-view";
import { useDraggable, useDroppable } from "@dnd-kit/core";
import { SortableContext, rectSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import Link from "next/link";
import {
  Fragment,
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

import { BulletIcon } from "@/components/ui/icons/Bullet";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { cx } from "@/styled-system/css";
import { type DragTreeVariantProps, dragTree } from "@/styled-system/recipes";

export type DragTreePosition = "top" | "in" | "bottom" | "only";
export type DragTreeDirection = "above" | "below";

type DragTreeNodeContext<TNode> = {
  node: TNode;
  parent: TNode | null;
  isRoot: boolean;
};

type DragTreeDividerContext<TNode> = DragTreeNodeContext<TNode> & {
  direction: DragTreeDirection;
};

export type DragTreeDragConfig<TNode> = {
  enabled: boolean;
  getNodeData: (context: DragTreeNodeContext<TNode>) => object;
  getDividerData: (context: DragTreeDividerContext<TNode>) => object;
};

export type DragTreeProps<TNode> = DragTreeVariantProps & {
  id: string;
  label: string;
  nodes: TNode[];
  getChildren: (node: TNode) => TNode[];
  getHref: (node: TNode) => string;
  getId: (node: TNode) => string;
  getLabel: (node: TNode) => string;
  defaultExpandedIds?: string[];
  drag?: DragTreeDragConfig<TNode>;
  getLabelClassName?: (
    context: DragTreeNodeContext<TNode>,
  ) => string | undefined;
  isSelected?: (node: TNode) => boolean;
  renderActions?: (context: DragTreeNodeContext<TNode>) => ReactNode;
  renderBefore?: (
    context: DragTreeNodeContext<TNode> & {
      index: number;
      siblings: TNode[];
    },
  ) => ReactNode;
};

type DragTreeMetadata = {
  treeId: string;
  kind: "node" | "divider";
  nodeId: string;
  direction?: DragTreeDirection;
};

type DragTreeEnvelope = {
  dragTree?: DragTreeMetadata;
};

export function DragTree<TNode>({
  id,
  label,
  nodes,
  getChildren,
  getHref,
  getId,
  getLabel,
  defaultExpandedIds = [],
  drag,
  getLabelClassName,
  isSelected,
  renderActions,
  renderBefore,
  variant,
}: DragTreeProps<TNode>) {
  const styles = dragTree({ variant });
  const collection = useMemo(
    () =>
      createTreeCollection<TNode>({
        nodeToValue: (node) => (node ? getId(node) : ""),
        nodeToString: (node) => (node ? getLabel(node) : ""),
        nodeToChildren: (node) => (node ? getChildren(node) : []),
        rootNode: { children: nodes } as TNode,
      }),
    [getChildren, getId, getLabel, nodes],
  );
  const [expandedValue, setExpandedValue] = useState(defaultExpandedIds);
  const previousDefaultExpandedIds = useRef(defaultExpandedIds);
  const selectedValue = useMemo(
    () =>
      findTreeNodeIds(
        nodes,
        getId,
        getChildren,
        (node) => isSelected?.(node) ?? false,
      ),
    [getChildren, getId, isSelected, nodes],
  );

  useEffect(() => {
    if (arraysEqual(previousDefaultExpandedIds.current, defaultExpandedIds)) {
      return;
    }

    previousDefaultExpandedIds.current = defaultExpandedIds;
    setExpandedValue((current) => [
      ...new Set([...current, ...defaultExpandedIds]),
    ]);
  }, [defaultExpandedIds]);

  const handleExpandedChange = useCallback(
    (details: TreeViewExpandedChangeDetails) => {
      setExpandedValue(details.expandedValue);
    },
    [],
  );

  const handleExpandNode = useCallback((nodeId: string) => {
    setExpandedValue((current) =>
      current.includes(nodeId) ? current : [...current, nodeId],
    );
  }, []);

  const handleToggleNode = useCallback((nodeId: string) => {
    setExpandedValue((current) =>
      current.includes(nodeId)
        ? current.filter((id) => id !== nodeId)
        : [...current, nodeId],
    );
  }, []);

  const rootNodes = nodes;

  return (
    <ArkTreeView.Root
      className={styles.root}
      collection={collection}
      expandedValue={expandedValue}
      onExpandedChange={handleExpandedChange}
      selectedValue={selectedValue}
    >
      <ArkTreeView.Label className={styles.label}>{label}</ArkTreeView.Label>
      <ArkTreeView.Tree className={styles.tree}>
        <SortableContext
          items={rootNodes.map(getId)}
          strategy={rectSortingStrategy}
        >
          {rootNodes.map((node, index) => {
            const context = {
              node,
              parent: null,
              isRoot: true,
              index,
              siblings: rootNodes,
            };

            return (
              <Fragment key={getId(node)}>
                {renderBefore?.(context)}
                <DragTreeNode
                  treeId={id}
                  allNodes={nodes}
                  node={node}
                  parent={null}
                  indexPath={[index]}
                  position={getDragTreePosition(rootNodes.length, index)}
                  getChildren={getChildren}
                  getHref={getHref}
                  getId={getId}
                  getLabel={getLabel}
                  drag={drag}
                  getLabelClassName={getLabelClassName}
                  handleExpandNode={handleExpandNode}
                  handleToggleNode={handleToggleNode}
                  isSelected={isSelected}
                  renderActions={renderActions}
                  styles={styles}
                />
              </Fragment>
            );
          })}
        </SortableContext>
      </ArkTreeView.Tree>
    </ArkTreeView.Root>
  );
}

type DragTreeNodeProps<TNode> = {
  treeId: string;
  allNodes: TNode[];
  node: TNode;
  parent: TNode | null;
  indexPath: number[];
  position: DragTreePosition;
  getChildren: (node: TNode) => TNode[];
  getHref: (node: TNode) => string;
  getId: (node: TNode) => string;
  getLabel: (node: TNode) => string;
  drag?: DragTreeDragConfig<TNode>;
  getLabelClassName?: (
    context: DragTreeNodeContext<TNode>,
  ) => string | undefined;
  handleExpandNode: (id: string) => void;
  handleToggleNode: (id: string) => void;
  isSelected?: (node: TNode) => boolean;
  renderActions?: (context: DragTreeNodeContext<TNode>) => ReactNode;
  styles: ReturnType<typeof dragTree>;
};

function DragTreeNode<TNode>({
  treeId,
  allNodes,
  node,
  parent,
  indexPath,
  position,
  getChildren,
  getHref,
  getId,
  getLabel,
  drag,
  getLabelClassName,
  handleExpandNode,
  handleToggleNode,
  isSelected,
  renderActions,
  styles,
}: DragTreeNodeProps<TNode>) {
  const nodeId = getId(node);
  const children = getChildren(node);
  const isRoot = parent === null;
  const context = { node, parent, isRoot };
  const nodeData = {
    ...drag?.getNodeData(context),
    dragTree: {
      treeId,
      kind: "node",
      nodeId,
    } satisfies DragTreeMetadata,
  };
  const {
    active,
    isDragging,
    listeners,
    over,
    setNodeRef: setDraggableNodeRef,
    transform,
  } = useDraggable({
    id: nodeId,
    data: nodeData,
    disabled: !drag?.enabled,
  });
  // The row contains links, disclosure buttons, and action controls. Let those
  // controls own keyboard input while retaining mouse and touch row dragging.
  const pointerListeners = listeners
    ? { ...listeners, onKeyDown: undefined }
    : undefined;
  const activeMetadata = getDragTreeMetadata(active?.data.current);
  const overMetadata = getDragTreeMetadata(over?.data.current);
  const draggedId =
    activeMetadata?.treeId === treeId && activeMetadata.kind === "node"
      ? activeMetadata.nodeId
      : null;
  const invalidTarget = draggedId
    ? draggedId === nodeId ||
      isDragTreeDescendant(allNodes, draggedId, nodeId, getId, getChildren)
    : false;
  const { setNodeRef: setDroppableNodeRef } = useDroppable({
    id: nodeId,
    data: nodeData,
    disabled: !drag?.enabled || invalidTarget,
    resizeObserverConfig: {
      updateMeasurementsFor: [],
    },
  });
  const isDraggingOver =
    !invalidTarget &&
    activeMetadata?.treeId === treeId &&
    activeMetadata.kind === "node" &&
    overMetadata?.treeId === treeId &&
    overMetadata.kind === "node" &&
    overMetadata.nodeId === nodeId &&
    !isDragging;
  const expandTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (!drag?.enabled || !isDraggingOver || children.length === 0) {
      if (expandTimeout.current) {
        clearTimeout(expandTimeout.current);
        expandTimeout.current = null;
      }
      return;
    }

    expandTimeout.current = setTimeout(() => {
      handleExpandNode(nodeId);
    }, 600);

    return () => {
      if (expandTimeout.current) {
        clearTimeout(expandTimeout.current);
        expandTimeout.current = null;
      }
    };
  }, [
    children.length,
    drag?.enabled,
    handleExpandNode,
    isDraggingOver,
    nodeId,
  ]);

  const showTopDivider =
    overMetadata?.treeId === treeId &&
    overMetadata.kind === "divider" &&
    overMetadata.nodeId === nodeId &&
    overMetadata.direction === "above";
  const showBottomDivider =
    overMetadata?.treeId === treeId &&
    overMetadata.kind === "divider" &&
    overMetadata.nodeId === nodeId &&
    overMetadata.direction === "below";
  const actions = renderActions?.(context);
  const selected = isSelected?.(node) ?? false;

  function preventNavigationWhileDragging(event: React.SyntheticEvent) {
    if (isDragging) {
      event.preventDefault();
      event.stopPropagation();
    }
  }

  return (
    <ArkTreeView.NodeProvider node={node} indexPath={indexPath}>
      <DragTreeDropIndicator
        active={showTopDivider}
        context={context}
        direction="above"
        disabled={invalidTarget}
        drag={drag}
        nodeId={nodeId}
        styles={styles}
        treeId={treeId}
      />

      <ArkTreeView.Branch
        ref={setDroppableNodeRef}
        className={styles.branch}
        data-position={position}
      >
        <div
          {...pointerListeners}
          ref={setDraggableNodeRef}
          className={cx("group", styles.row)}
          data-drag-over={isDraggingOver}
          data-selected={selected}
          style={{
            opacity: isDragging ? 0.5 : 1,
            pointerEvents: isDragging ? "none" : undefined,
            transform: CSS.Translate.toString(transform),
          }}
        >
          {children.length > 0 ? (
            <button
              aria-label={`Toggle ${getLabel(node)}`}
              className={styles.branchControl}
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                handleToggleNode(nodeId);
              }}
              onKeyDown={(event) => {
                event.stopPropagation();
                if (event.key === "Enter" || event.key === " ") {
                  event.preventDefault();
                  handleToggleNode(nodeId);
                }
              }}
              onPointerDown={(event) => event.stopPropagation()}
              onTouchStart={(event) => event.stopPropagation()}
            >
              <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
                <ChevronRightIcon />
              </ArkTreeView.BranchIndicator>
            </button>
          ) : (
            <span aria-hidden="true" className={styles.branchControl}>
              <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
                <BulletIcon />
              </ArkTreeView.BranchIndicator>
            </span>
          )}

          <ArkTreeView.BranchText
            asChild
            className={cx(styles.branchText, getLabelClassName?.(context))}
          >
            <Link
              aria-current={selected ? "page" : undefined}
              className={styles.link}
              href={getHref(node)}
              onClick={preventNavigationWhileDragging}
              onTouchStart={preventNavigationWhileDragging}
            >
              {getLabel(node)}
            </Link>
          </ArkTreeView.BranchText>

          {actions && (
            <div
              className={styles.actions}
              onKeyDown={(event) => event.stopPropagation()}
            >
              {actions}
            </div>
          )}
        </div>

        <ArkTreeView.BranchContent className={styles.branchContent}>
          <SortableContext
            items={children.map(getId)}
            strategy={rectSortingStrategy}
          >
            {children.map((child, index) => (
              <DragTreeNode
                key={getId(child)}
                treeId={treeId}
                allNodes={allNodes}
                node={child}
                parent={node}
                indexPath={[...indexPath, index]}
                position={getDragTreePosition(children.length, index)}
                getChildren={getChildren}
                getHref={getHref}
                getId={getId}
                getLabel={getLabel}
                drag={drag}
                getLabelClassName={getLabelClassName}
                handleExpandNode={handleExpandNode}
                handleToggleNode={handleToggleNode}
                isSelected={isSelected}
                renderActions={renderActions}
                styles={styles}
              />
            ))}
          </SortableContext>
        </ArkTreeView.BranchContent>
      </ArkTreeView.Branch>

      {(position === "bottom" || position === "only") && (
        <DragTreeDropIndicator
          active={showBottomDivider}
          context={context}
          direction="below"
          disabled={invalidTarget}
          drag={drag}
          nodeId={nodeId}
          styles={styles}
          treeId={treeId}
        />
      )}
    </ArkTreeView.NodeProvider>
  );
}

type DragTreeDropIndicatorProps<TNode> = {
  active: boolean;
  context: DragTreeNodeContext<TNode>;
  direction: DragTreeDirection;
  disabled: boolean;
  drag?: DragTreeDragConfig<TNode>;
  nodeId: string;
  styles: ReturnType<typeof dragTree>;
  treeId: string;
};

function DragTreeDropIndicator<TNode>({
  active,
  context,
  direction,
  disabled,
  drag,
  nodeId,
  styles,
  treeId,
}: DragTreeDropIndicatorProps<TNode>) {
  const data = {
    ...drag?.getDividerData({ ...context, direction }),
    dragTree: {
      treeId,
      kind: "divider",
      nodeId,
      direction,
    } satisfies DragTreeMetadata,
  };
  const { setNodeRef } = useDroppable({
    id: `${nodeId}_${direction}`,
    data,
    disabled: !drag?.enabled || disabled,
    resizeObserverConfig: {
      updateMeasurementsFor: [],
    },
  });

  if (!drag?.enabled) {
    return null;
  }

  return (
    <div
      className={styles.dropIndicator}
      data-direction={direction}
      data-node-id={nodeId}
      data-part="drop-indicator"
    >
      <div
        ref={setNodeRef}
        className={styles.dropIndicatorLine}
        data-active={active}
      />
    </div>
  );
}

function getDragTreeMetadata(
  data: Record<string, unknown> | undefined,
): DragTreeMetadata | null {
  const metadata = (data as DragTreeEnvelope | undefined)?.dragTree;
  if (!metadata) {
    return null;
  }

  return metadata;
}

export function getDragTreePosition(
  itemCount: number,
  index: number,
): DragTreePosition {
  if (itemCount === 1) {
    return "only";
  }
  if (index === 0) {
    return "top";
  }
  if (index === itemCount - 1) {
    return "bottom";
  }
  return "in";
}

export function isDragTreeDescendant<TNode>(
  nodes: TNode[],
  ancestorId: string,
  descendantId: string,
  getId: (node: TNode) => string,
  getChildren: (node: TNode) => TNode[],
): boolean {
  const ancestor = findTreeNode(nodes, ancestorId, getId, getChildren);
  if (!ancestor) {
    return false;
  }

  return (
    findTreeNode(getChildren(ancestor), descendantId, getId, getChildren) !==
    null
  );
}

function findTreeNode<TNode>(
  nodes: TNode[],
  id: string,
  getId: (node: TNode) => string,
  getChildren: (node: TNode) => TNode[],
): TNode | null {
  for (const node of nodes) {
    if (getId(node) === id) {
      return node;
    }
    const descendant = findTreeNode(getChildren(node), id, getId, getChildren);
    if (descendant) {
      return descendant;
    }
  }
  return null;
}

function findTreeNodeIds<TNode>(
  nodes: TNode[],
  getId: (node: TNode) => string,
  getChildren: (node: TNode) => TNode[],
  predicate: (node: TNode) => boolean,
): string[] {
  const ids: string[] = [];
  const visit = (node: TNode) => {
    if (predicate(node)) {
      ids.push(getId(node));
    }
    getChildren(node).forEach(visit);
  };
  nodes.forEach(visit);
  return ids;
}

function arraysEqual(left: string[], right: string[]) {
  return (
    left.length === right.length &&
    left.every((value, index) => value === right[index])
  );
}
