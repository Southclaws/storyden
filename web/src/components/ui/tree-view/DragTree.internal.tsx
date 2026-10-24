"use client";

import {
  TreeView as ArkTreeView,
  createTreeCollection,
  useTreeViewContext,
} from "@ark-ui/react/tree-view";
import { useDraggable, useDroppable } from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import Link from "next/link";
import {
  type CSSProperties,
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
import { Tooltip } from "@/components/ui/tooltip/Tooltip";
import { cx } from "@/styled-system/css";
import { type DragTreeVariantProps, dragTree } from "@/styled-system/recipes";

import { DragTreeMovePreview } from "./DragTreeMovePreview";
import { type DragTreeMoveConfig, useDragTreeMove } from "./useDragTreeMove";

export type { DragTreeMove, DragTreeMoveConfig } from "./useDragTreeMove";

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

type DragTreeActionContext<TNode> = DragTreeNodeContext<TNode> & {
  /** Only the focused node's actions belong in the Tab order. */
  tabIndex: 0 | -1;
  startMove?: () => void;
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
  move?: DragTreeMoveConfig<TNode>;
  getLabelClassName?: (
    context: DragTreeNodeContext<TNode>,
  ) => string | undefined;
  isSelected?: (node: TNode) => boolean;
  renderActions?: (context: DragTreeActionContext<TNode>) => ReactNode;
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
  move,
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
        nodeToValue: getId,
        nodeToString: getLabel,
        nodeToChildren: getChildren,
        rootNode: { children: nodes } as TNode,
      }),
    [getChildren, getId, getLabel, nodes],
  );
  const selectedValue = collection.getValues().filter((value) => {
    const node = collection.findNode(value);
    return node && isSelected?.(node);
  });
  const [expandedValue, setExpandedValue] = useState(defaultExpandedIds);
  const previousDefaultExpandedIds = useRef(defaultExpandedIds);

  useEffect(() => {
    if (arraysEqual(previousDefaultExpandedIds.current, defaultExpandedIds)) {
      return;
    }

    previousDefaultExpandedIds.current = defaultExpandedIds;
    setExpandedValue((current) => [
      ...new Set([...current, ...defaultExpandedIds]),
    ]);
  }, [defaultExpandedIds]);

  const handleExpandNode = useCallback((nodeId: string) => {
    setExpandedValue((current) =>
      current.includes(nodeId) ? current : [...current, nodeId],
    );
  }, []);
  const rootRef = useRef<HTMLDivElement>(null);
  const keyboardMove = useDragTreeMove({
    rootRef,
    collection,
    config: move,
    enabled: !!drag?.enabled,
    expand: handleExpandNode,
  });

  const rootNodes = nodes;

  if (nodes.length === 0) return null;

  return (
    <ArkTreeView.Root
      ref={rootRef}
      className={styles.root}
      collection={collection}
      expandedValue={expandedValue}
      onExpandedChange={(details) => setExpandedValue(details.expandedValue)}
      selectedValue={selectedValue}
      defaultFocusedValue={selectedValue[0] ?? (nodes[0] && getId(nodes[0]))}
      expandOnClick={false}
      onKeyDownCapture={keyboardMove.onKeyDown}
      onPointerDownCapture={() => {
        keyboardMove.setKeyboard(false);
        keyboardMove.cancel(false);
      }}
      onFocusCapture={(event) => {
        const element = event.target as HTMLElement;
        if (element.dataset["treeNodeId"]) {
          keyboardMove.setFocusedId(element.dataset["treeNodeId"]);
          keyboardMove.setKeyboard(element.matches(":focus-visible"));
        } else keyboardMove.setKeyboard(false);
      }}
      onBlurCapture={(event) => {
        if (
          event.relatedTarget &&
          event.currentTarget.contains(event.target as Node) &&
          !event.currentTarget.contains(event.relatedTarget as Node | null)
        ) {
          keyboardMove.setKeyboard(false);
          keyboardMove.cancel(false);
        }
      }}
    >
      <ArkTreeView.Label asChild className={styles.label}>
        <p>{label}</p>
      </ArkTreeView.Label>
      <ArkTreeView.Tree className={styles.tree}>
        {rootNodes.map((node, index) => {
          const context = {
            node,
            parent: null,
            isRoot: true,
            index,
            siblings: rootNodes,
          };
          const before = renderBefore?.(context);

          return (
            <Fragment key={getId(node)}>
              {before}
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
                keyboardMove={keyboardMove}
                moveEnabled={!!move && !!drag?.enabled}
                getLabelClassName={getLabelClassName}
                handleExpandNode={handleExpandNode}
                isSelected={isSelected}
                renderActions={renderActions}
                styles={styles}
              />
            </Fragment>
          );
        })}
      </ArkTreeView.Tree>
      <p id={`${id}-move-instructions`} className={styles.instructions}>
        Alt+Shift+M picks up a node. Use arrows and typeahead to browse
        destinations. Space cycles before, inside, and after. Enter places;
        Escape cancels. Tab cancels and leaves move mode.
      </p>
      <div
        role="status"
        aria-live="polite"
        aria-atomic="true"
        className={styles.instructions}
      >
        {keyboardMove.status}
      </div>
      {keyboardMove.error && (
        <p role="alert" className={styles.moveError}>
          {keyboardMove.error}
        </p>
      )}
      {keyboardMove.source && (
        <DragTreeMovePreview
          rootRef={rootRef}
          focusedId={keyboardMove.focusedId}
          label={getLabel(keyboardMove.source)}
          preview={keyboardMove.preview}
          pending={keyboardMove.pending}
          styles={styles}
        />
      )}
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
  keyboardMove: ReturnType<typeof useDragTreeMove<TNode>>;
  moveEnabled: boolean;
  getLabelClassName?: (
    context: DragTreeNodeContext<TNode>,
  ) => string | undefined;
  handleExpandNode: (id: string) => void;
  isSelected?: (node: TNode) => boolean;
  renderActions?: (context: DragTreeActionContext<TNode>) => ReactNode;
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
  keyboardMove,
  moveEnabled,
  getLabelClassName,
  handleExpandNode,
  isSelected,
  renderActions,
  styles,
}: DragTreeNodeProps<TNode>) {
  const nodeId = getId(node);
  const children = getChildren(node);
  const linkRef = useRef<HTMLAnchorElement>(null);
  const tree = useTreeViewContext();
  const nodeState = tree.getNodeState({ node, indexPath });
  const { expanded, focused, depth } = nodeState;
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
    disabled: !drag?.enabled || !!keyboardMove.sourceId,
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
    disabled: !drag?.enabled || invalidTarget || !!keyboardMove.sourceId,
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

  const keyboardTarget =
    !!keyboardMove.sourceId &&
    keyboardMove.focusedId === nodeId &&
    !keyboardMove.invalid;
  const showTopDivider =
    (keyboardTarget && keyboardMove.placement === "above") ||
    (overMetadata?.treeId === treeId &&
      overMetadata.kind === "divider" &&
      overMetadata.nodeId === nodeId &&
      overMetadata.direction === "above");
  const showBottomDivider =
    (keyboardTarget && keyboardMove.placement === "below") ||
    (overMetadata?.treeId === treeId &&
      overMetadata.kind === "divider" &&
      overMetadata.nodeId === nodeId &&
      overMetadata.direction === "below");
  const actions = renderActions?.({
    ...context,
    tabIndex: focused && !keyboardMove.sourceId ? 0 : -1,
    startMove: moveEnabled ? () => keyboardMove.start(nodeId) : undefined,
  });
  const selected = isSelected?.(node) ?? false;

  function preventNavigationWhileDragging(event: React.SyntheticEvent) {
    if (isDragging || keyboardMove.sourceId) {
      event.preventDefault();
      event.stopPropagation();
    }
  }

  const Branch = children.length > 0 ? ArkTreeView.Branch : "div";
  // Ark's focus and keyboard props belong on the destination link. The outer
  // tree item owns the hierarchy and selection state, including row actions.
  const linkProps =
    children.length > 0
      ? tree.getBranchControlProps({ node, indexPath })
      : tree.getItemProps({ node, indexPath });
  const link = (
    <Link
      {...linkProps}
      ref={linkRef}
      role="link"
      aria-level={undefined}
      aria-selected={undefined}
      data-tree-node-id={nodeId}
      aria-keyshortcuts={moveEnabled ? "Alt+Shift+M" : undefined}
      aria-describedby={moveEnabled ? `${treeId}-move-instructions` : undefined}
      aria-current={selected ? "page" : undefined}
      aria-expanded={children.length > 0 ? expanded : undefined}
      className={cx(
        styles.branchText,
        styles.link,
        getLabelClassName?.(context),
      )}
      href={getHref(node)}
      prefetch={false}
      onClick={(event) => {
        preventNavigationWhileDragging(event);
        if (!event.defaultPrevented) {
          linkProps.onClick?.(event);
        }
      }}
      onTouchStart={preventNavigationWhileDragging}
    >
      {getLabel(node)}
    </Link>
  );

  const row = (
    <div
      {...(keyboardMove.sourceId ? undefined : pointerListeners)}
      ref={setDraggableNodeRef}
      className={cx("group", styles.row)}
      data-drag-over={isDraggingOver}
      data-selected={selected}
      data-depth={depth}
      data-carried={keyboardMove.sourceId === nodeId}
      data-move-target={keyboardTarget ? keyboardMove.placement : undefined}
      style={{
        opacity: isDragging ? 0.5 : 1,
        pointerEvents: isDragging ? "none" : undefined,
        transform: CSS.Translate.toString(transform),
      }}
    >
      {children.length > 0 ? (
        <ArkTreeView.BranchTrigger
          asChild
          aria-expanded={expanded}
          className={styles.branchControl}
          tabIndex={-1}
          onPointerDown={(event) => event.stopPropagation()}
          onMouseDown={(event) => event.stopPropagation()}
          onTouchStart={(event) => event.stopPropagation()}
        >
          <button type="button" aria-label={`Toggle ${getLabel(node)}`}>
            <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
              <ChevronRightIcon />
            </ArkTreeView.BranchIndicator>
          </button>
        </ArkTreeView.BranchTrigger>
      ) : (
        <span aria-hidden="true" className={styles.branchControl}>
          <span className={styles.branchIndicator}>
            <BulletIcon />
          </span>
        </span>
      )}

      {link}

      {actions && (
        // oxlint-disable-next-line jsx-a11y/no-noninteractive-element-interactions -- Delegate child control events without letting them reach the tree's keyboard or drag handlers.
        <div
          role="group"
          aria-label={`${getLabel(node)} actions`}
          className={styles.actions}
          inert={!!keyboardMove.sourceId}
          onPointerDown={(event) => event.stopPropagation()}
          onMouseDown={(event) => event.stopPropagation()}
          onTouchStart={(event) => event.stopPropagation()}
          onKeyDown={(event) => {
            event.stopPropagation();
            if (
              event.key === "Escape" &&
              !event.defaultPrevented &&
              event.currentTarget.contains(event.target as Node)
            ) {
              event.preventDefault();
              linkRef.current?.focus();
            }
          }}
        >
          {actions}
        </div>
      )}
      {moveEnabled && (
        <Tooltip
          open={
            keyboardMove.hintId === nodeId && focused && !keyboardMove.sourceId
          }
          positioning={{
            placement: "right-start",
            gutter: 12,
            getAnchorRect: () =>
              linkRef.current?.parentElement?.getBoundingClientRect() ?? null,
          }}
          content={
            <span>Move or reorder: Alt+Shift+M (Option+Shift+M on Mac)</span>
          }
          contentProps={{
            _open: { animation: "[none]" },
            _closed: { display: "none" },
          }}
        >
          <span aria-hidden="true" className={styles.instructions} />
        </Tooltip>
      )}
    </div>
  );

  return (
    <ArkTreeView.NodeProvider node={node} indexPath={indexPath}>
      <DragTreeDropIndicator
        active={showTopDivider}
        context={context}
        direction="above"
        disabled={invalidTarget || !!keyboardMove.sourceId}
        drag={drag}
        nodeId={nodeId}
        styles={styles}
        treeId={treeId}
      />

      <Branch
        ref={setDroppableNodeRef}
        className={styles.branch}
        data-position={position}
        role="treeitem"
        aria-labelledby={nodeState.id}
        aria-level={depth}
        aria-selected={nodeState.selected}
        style={{ "--depth": depth } as CSSProperties}
      >
        {row}

        {children.length > 0 && (
          <ArkTreeView.BranchContent className={styles.branchContent}>
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
                keyboardMove={keyboardMove}
                moveEnabled={moveEnabled}
                getLabelClassName={getLabelClassName}
                handleExpandNode={handleExpandNode}
                isSelected={isSelected}
                renderActions={renderActions}
                styles={styles}
              />
            ))}
          </ArkTreeView.BranchContent>
        )}
      </Branch>

      {(position === "bottom" || position === "only" || showBottomDivider) && (
        <DragTreeDropIndicator
          active={showBottomDivider}
          context={context}
          direction="below"
          disabled={invalidTarget || !!keyboardMove.sourceId}
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

function arraysEqual(left: string[], right: string[]) {
  return (
    left.length === right.length &&
    left.every((value, index) => value === right[index])
  );
}
