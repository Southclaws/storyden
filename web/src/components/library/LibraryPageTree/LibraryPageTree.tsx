import {
  TreeView as ArkTreeView,
  createTreeCollection,
} from "@ark-ui/react/tree-view";
import { useDraggable, useDroppable } from "@dnd-kit/core";
import { SortableContext, rectSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import Link from "next/link";
import { ExpandedChangeDetails } from "node_modules/@ark-ui/react/dist/components/tree-view/tree-view";
import {
  CSSProperties,
  Fragment,
  JSX,
  useEffect,
  useRef,
  useState,
} from "react";

import { Identifier, NodeWithChildren, Visibility } from "@/api/openapi-schema";
import { CreatePageAction } from "@/components/library/CreatePage";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { ReviewIcon, UnlistedIcon } from "@/components/ui/icons/Visibility";
import {
  DragItemData,
  DragItemDivider,
  DragItemNode,
} from "@/lib/dragdrop/provider";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";

import { LibraryPageMenu } from "../LibraryPageMenu/LibraryPageMenu";

export interface LibraryPageTreeProps {
  nodes: NodeWithChildren[];
  currentNode: string | undefined;
  canManageLibrary: boolean;
}

const visibilitySortKey: Record<Visibility, number> = {
  [Visibility.published]: 0,
  [Visibility.review]: 1,
  [Visibility.draft]: 2,
  [Visibility.unlisted]: 3,
};

const visibilityLabels: Record<Visibility, string> = {
  [Visibility.published]: "Published",
  [Visibility.review]: "In review",
  [Visibility.draft]: "Drafts",
  [Visibility.unlisted]: "Unlisted",
};

const visibilityIcons: Record<Visibility, JSX.Element> = {
  [Visibility.published]: <></>,
  [Visibility.review]: <ReviewIcon />,
  [Visibility.draft]: <DraftIcon />,
  [Visibility.unlisted]: <UnlistedIcon />,
};

export type PositionInList = "top" | "in" | "bottom" | "only";

export function getPositionInList(
  numberOfNodes: number,
  index: number,
): PositionInList {
  if (numberOfNodes === 1) {
    return "only";
  }

  if (index === 0) {
    return "top";
  }

  if (index === numberOfNodes - 1) {
    return "bottom";
  }

  return "in";
}

export const LibraryPageTree = (props: LibraryPageTreeProps) => {
  const { nodes, currentNode } = props;

  const styles = treeView();

  const defaultExpandedValue: string[] = [];

  const findCurrentNode = (node: NodeWithChildren) => {
    if (node.slug === currentNode) {
      defaultExpandedValue.push(node.id);
      return true;
    }

    if (node.children) {
      for (const child of node.children) {
        if (findCurrentNode(child)) {
          defaultExpandedValue.push(node.id);
          return true;
        }
      }
    }

    return false;
  };

  nodes.forEach(findCurrentNode);

  const sortedByVisibility = nodes.sort((a, b) => {
    return visibilitySortKey[a.visibility] - visibilitySortKey[b.visibility];
  });

  const collection = createTreeCollection<NodeWithChildren>({
    // NOTE: Ark bug where sometimes these functions receive an undefined value.
    nodeToValue: (n?: NodeWithChildren) => {
      return n?.id ?? "";
    },
    nodeToString: (n?: NodeWithChildren) => {
      return n?.slug ?? "";
    },
    nodeToChildren: (n?: NodeWithChildren) => {
      return n?.children ?? [];
    },
    rootNode: {
      children: sortedByVisibility,
    } as NodeWithChildren,
  });

  const rootNodes = collection.rootNode.children;

  const [expandedValue, setExpandedValue] = useState(defaultExpandedValue);

  function handleExpandedChange(e: ExpandedChangeDetails) {
    setExpandedValue(e.expandedValue);
  }

  function handleExpandNode(id: string) {
    setExpandedValue((prev) => {
      return [...prev, id];
    });
  }

  return (
    <ArkTreeView.Root
      className={styles.root}
      collection={collection}
      defaultExpandedValue={defaultExpandedValue}
      expandedValue={expandedValue}
      onExpandedChange={handleExpandedChange}
    >
      <ArkTreeView.Tree className={cx(styles.tree)}>
        <SortableContext items={rootNodes.map((child) => child.id)}>
          {rootNodes.map((node, index) => {
            const previous = index > 0 ? rootNodes[index - 1] : null;

            const sameVisibilityAsPrevious = previous
              ? previous.visibility === node.visibility
              : true;

            const dividerLabel = visibilityLabels[node.visibility];
            const dividerIcon = visibilityIcons[node.visibility];

            // We only show dividers on the root list, as this is the only list that's
            // sorted by visibility.
            const showDivider = !sameVisibilityAsPrevious;

            return (
              <Fragment key={node.id}>
                {showDivider && (
                  <HStack mb="1">
                    <NavigationHeader href="/drafts">
                      <HStack>
                        {dividerIcon}
                        {dividerLabel}
                      </HStack>
                    </NavigationHeader>
                  </HStack>
                )}

                <TreeNode
                  fullTree={nodes}
                  currentNode={currentNode}
                  parentID={null}
                  node={node}
                  indexPath={[]}
                  isRoot={true}
                  styles={styles}
                  positionInList={getPositionInList(rootNodes.length, index)}
                  handleExpandNode={handleExpandNode}
                  canManageLibrary={props.canManageLibrary}
                />
              </Fragment>
            );
          })}
        </SortableContext>
      </ArkTreeView.Tree>
    </ArkTreeView.Root>
  );
};

type TreeNodeProps = {
  fullTree: NodeWithChildren[];
  currentNode: string | undefined;
  parentID: Identifier | null;
  node: NodeWithChildren;
  styles: any;
  isRoot: boolean;
  indexPath: number[];
  positionInList: PositionInList;
  handleExpandNode: (id: string) => void;
  canManageLibrary: boolean;
};

const linkStyles = css({
  // disable ios preview
  touchAction: "none",
  userSelect: "none",
  WebkitTouchCallout: "none",
});

function TreeNode({
  fullTree,
  currentNode,
  parentID,
  styles,
  node,
  isRoot,
  indexPath,
  positionInList,
  handleExpandNode,
  canManageLibrary,
}: TreeNodeProps) {
  const {
    attributes,
    listeners,
    setNodeRef: setDraggableNodeRef,
    transform,
    isDragging,
    active,
    over,
  } = useDraggable({
    disabled: !canManageLibrary,
    id: node.id,
    data: {
      type: "node",
      node,
      parentID: parentID,
      context: "sidebar",
    } satisfies DragItemNode,
  });

  const dragged = active?.data.current as DragItemData | undefined;
  const draggedID = dragged?.type === "node" ? dragged.node.id : undefined;

  const isDescendantOfDraggedNode = draggedID
    ? isDescendant(fullTree, draggedID, node.id)
    : false;

  const { setNodeRef: setDroppableNodeRef } = useDroppable({
    disabled: isDescendantOfDraggedNode || !canManageLibrary,
    id: node.id,
    data: {
      type: "node",
      node,
      parentID: parentID,
      context: "sidebar",
    } satisfies DragItemNode,
    resizeObserverConfig: {
      updateMeasurementsFor: [],
    },
  });

  const overItem = over?.data.current as DragItemData | undefined;

  const isDraggingOver =
    !isDescendantOfDraggedNode &&
    overItem?.type === "node" &&
    dragged?.type === "node" &&
    overItem?.node.id === node.id &&
    overItem.context === "sidebar" &&
    !isDragging;

  // handle drag-over to expand
  const expandTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  useEffect(() => {
    if (isDraggingOver) {
      expandTimeout.current = setTimeout(() => {
        handleExpandNode(node.id);
      }, 600);
    } else {
      if (expandTimeout.current) {
        clearTimeout(expandTimeout.current);
        expandTimeout.current = null;
      }
    }

    return () => {
      if (expandTimeout.current) {
        clearTimeout(expandTimeout.current);
        expandTimeout.current = null;
      }
    };
  }, [isDraggingOver, node.id, handleExpandNode]);

  function handleLinkClick(e: React.MouseEvent) {
    if (isDragging) {
      e.preventDefault();
      e.stopPropagation();
    }
  }

  function handleLinkTouchStart(e: React.TouchEvent) {
    if (isDragging) {
      e.preventDefault();
    }
  }

  const branchControlDragStyles: CSSProperties = {
    transform: CSS.Translate.toString(transform),
    opacity: isDragging ? 0.5 : 1,
    ...(isDragging
      ? {
          pointerEvents: "none",
        }
      : {
          pointerEvents: "unset",
        }),
  };

  const showTopDivider =
    overItem?.type === "divider" &&
    overItem?.siblingNode.id === node.id &&
    overItem?.direction === "above";
  const showBottomDivider =
    overItem?.type === "divider" &&
    overItem?.siblingNode.id === node.id &&
    overItem?.direction === "below";

  const isPublished = node.visibility === Visibility.published;
  const isHighlighted = node.slug === currentNode;

  const label = isPublished ? node.name : `${node.name}`;

  const branchColourPalette = visibilityColour(node.visibility);

  const visibilityStyles = isRoot
    ? "" // Don't show the visibility state styles for root items, is cluttered.
    : css({
        paddingX: "0.5",
        borderRadius: "sm",
        colorPalette: branchColourPalette,
        borderWidth: node.visibility === Visibility.published ? "none" : "thin",
        borderColor:
          node.visibility === Visibility.published
            ? "transparent"
            : "colorPalette.6",
        borderStyle:
          node.visibility === Visibility.published ? "solid" : "dashed",
      });

  const draggingOverStyles = isDraggingOver
    ? css({
        borderRadius: "md",
        colorPalette: branchColourPalette,
        outlineWidth: "thin",
        outlineStyle: "dashed",
        outlineColor: "colorPalette.6",
        outlineOffset: "-0.5",
      })
    : "";

  const highlightStyles = css({
    background: isHighlighted ? "bg.selected" : undefined,
  });

  return (
    <ArkTreeView.NodeProvider key={node.id} node={node} indexPath={indexPath}>
      <DropIndicator
        node={node}
        direction="above"
        active={showTopDivider}
        positionInList={positionInList}
        parentID={parentID}
      />

      <ArkTreeView.Branch
        ref={setDroppableNodeRef}
        className={cx(styles.branch)}
        data-visibility={node.visibility}
        data-position={positionInList}
      >
        <ArkTreeView.BranchControl
          className={cx(
            "group",
            styles.branchControl,
            highlightStyles,
            draggingOverStyles,
          )}
          ref={setDraggableNodeRef}
          style={branchControlDragStyles}
          {...attributes}
          {...listeners}
        >
          <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
            {node.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
          </ArkTreeView.BranchIndicator>

          <ArkTreeView.BranchText
            asChild
            className={cx(styles.branchText, visibilityStyles)}
          >
            <Link
              className={linkStyles}
              onClick={handleLinkClick}
              onTouchStart={handleLinkTouchStart}
              href={`/l/${node.slug}`}
            >
              {label}
            </Link>
          </ArkTreeView.BranchText>

          <HStack
            opacity={{
              base: "0",
              _groupHover: "full",
            }}
            gap="1"
            minW="min"
            flexShrink="0"
          >
            {canManageLibrary && (
              <CreatePageAction
                variant="ghost"
                hideLabel
                parentSlug={node.slug}
              />
            )}
            <LibraryPageMenu variant="ghost" node={node} />
          </HStack>
        </ArkTreeView.BranchControl>

        <ArkTreeView.BranchContent className={styles.branchContent}>
          <SortableContext
            items={node.children.map((child) => child.id)}
            strategy={rectSortingStrategy}
          >
            {node.children.map((child, index) => {
              return (
                <TreeNode
                  key={child.id}
                  fullTree={fullTree}
                  currentNode={currentNode}
                  parentID={node.id}
                  node={child}
                  indexPath={[...indexPath, index]}
                  isRoot={false}
                  styles={styles}
                  positionInList={getPositionInList(
                    node.children.length,
                    index,
                  )}
                  handleExpandNode={handleExpandNode}
                  canManageLibrary={canManageLibrary}
                />
              );
            })}
          </SortableContext>
        </ArkTreeView.BranchContent>
      </ArkTreeView.Branch>
      {positionInList === "bottom" ||
        (positionInList === "only" && (
          <DropIndicator
            node={node}
            direction="below"
            active={showBottomDivider}
            positionInList={positionInList}
            parentID={parentID}
          />
        ))}
    </ArkTreeView.NodeProvider>
  );
}

function DropIndicator({
  node,
  direction,
  active,
  positionInList,
  parentID,
}: {
  parentID: Identifier | null;
  node: NodeWithChildren;
  direction: "above" | "below";
  active: boolean;
  positionInList: PositionInList;
}) {
  const { setNodeRef } = useDroppable({
    id: `${node.id}_${direction}`,
    data: {
      type: "divider",
      direction,
      siblingNode: node,
      parentID: parentID,
      context: "sidebar",
    } satisfies DragItemDivider,
    resizeObserverConfig: {
      updateMeasurementsFor: [],
    },
  });

  return (
    <div
      data-divider-node-id={node.id}
      data-divider-direction={direction}
      data-divider-active={active}
      data-divider-position={positionInList}
      style={{
        position: "relative",
        height: "1px",
        marginInlineStart: "calc(((var(--depth)) * 22px) + 22px)",
      }}
    >
      <div
        ref={setNodeRef}
        style={{
          position: "absolute",
          top: "-1px",
          left: 0,
          right: 0,
          height: "3px",
          background: active ? "var(--colors-bg-muted)" : "transparent",
          opacity: active ? 1 : 0,
          transition: "opacity 0.2s",
          pointerEvents: "none",
        }}
      />
    </div>
  );
}

export function isDescendant(
  nodes: NodeWithChildren[],
  ancestorId: string,
  descendantId: string,
): boolean {
  function dfs(node: NodeWithChildren): boolean {
    if (node.id === descendantId) {
      return true;
    }
    return node.children?.some(dfs) ?? false;
  }

  const ancestorNode = nodes.find((n) => n.id === ancestorId);
  if (!ancestorNode) return false;
  return dfs(ancestorNode);
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
