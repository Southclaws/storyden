import type { Assign } from "@ark-ui/react";
import {
  TreeView as ArkTreeView,
  type TreeViewRootProps,
} from "@ark-ui/react/tree-view";
import {
  DndContext,
  DragOverlay,
  KeyboardSensor,
  Over,
  PointerSensor,
  closestCenter,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { keyBy } from "lodash";
import Link from "next/link";
import { MouseEvent, PropsWithChildren, forwardRef, useState } from "react";
import { createPortal } from "react-dom";

import { NodeWithChildren, Visibility } from "@/api/openapi-schema";
import { CreatePageAction } from "@/components/site/Navigation/Actions/CreatePage";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx } from "@/styled-system/css";
import { Box, HStack, splitCssProps } from "@/styled-system/jsx";
import { type TreeViewVariantProps, treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";
import type { JsxStyleProps } from "@/styled-system/types";

import { LibraryPageMenu } from "../LibraryPageMenu/LibraryPageMenu";

export interface TreeViewData {
  label: string;
  children: NodeWithChildren[];
}

export interface TreeViewProps
  extends Assign<JsxStyleProps, TreeViewRootProps>,
    TreeViewVariantProps {
  data: TreeViewData;
  currentNode: string | undefined;
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

export const LibraryPageTree = forwardRef<HTMLDivElement, TreeViewProps>(
  (props, ref) => {
    const [cssProps, localProps] = splitCssProps(props);
    const { data, currentNode, className, ...rootProps } = localProps;

    const styles = treeView();

    const defaultExpandedValue: string[] = [];

    const rootNodeMap = keyBy(data.children, (child) => child.id);

    // recursively find currentNode in data and add each parent to defaultExpandedValue
    const findCurrentNode = (node: NodeWithChildren) => {
      if (node.slug === currentNode) {
        defaultExpandedValue.push(node.slug);
        return true;
      }

      if (node.children) {
        for (const child of node.children) {
          if (findCurrentNode(child)) {
            defaultExpandedValue.push(node.slug);
            return true;
          }
        }
      }

      return false;
    };

    data.children.forEach(findCurrentNode);

    const renderChild = (child: NodeWithChildren, index: number) => {
      const previous = index > 0 ? data.children[index - 1] : null;

      const sameVisibilityAsPrevious = previous
        ? previous.visibility === child.visibility
        : true;

      const isRoot = Boolean(rootNodeMap[child.id]);

      const dividerLabel = visibilityLabels[child.visibility];

      // We only show dividers on the root list, as this is the only list that's
      // sorted by visibility.
      const showDivider = isRoot && !sameVisibilityAsPrevious;

      const isHighlighted = child.slug === currentNode;

      return (
        <DroppableBranch
          key={child.id}
          child={child}
          className={styles.branch}
          render={({ isOver, over }) => {
            return (
              <>
                {showDivider && (
                  <HStack mb="1">
                    <NavigationHeader href="/drafts">
                      {dividerLabel}
                    </NavigationHeader>
                  </HStack>
                )}

                <TreeBranch
                  styles={styles}
                  child={child}
                  isHovered={isOver}
                  isHighlighted={isHighlighted}
                  isRoot={isRoot}
                />

                <ArkTreeView.BranchContent className={styles.branchContent}>
                  <SortableContext
                    items={child.children}
                    // strategy={verticalListSortingStrategy}
                  >
                    {child.children?.map((child, i) =>
                      child.children ? (
                        renderChild(child, i)
                      ) : (
                        <TreeItem
                          key={child.id}
                          styles={styles}
                          child={child}
                          isHighlighted={isHighlighted}
                          isRoot={isRoot}
                        />
                      ),
                    )}
                  </SortableContext>
                </ArkTreeView.BranchContent>
              </>
            );
          }}
        />
      );
    };

    const sortedByVisibility = data.children.sort((a, b) => {
      return visibilitySortKey[a.visibility] - visibilitySortKey[b.visibility];
    });

    const sensors = useSensors(
      useSensor(PointerSensor, {
        activationConstraint: {
          distance: 4,
        },
      }),
      useSensor(KeyboardSensor, {
        coordinateGetter: sortableKeyboardCoordinates,
      }),
    );

    function handleDragEnd(e) {
      console.log(e);
    }

    return (
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
      >
        <ArkTreeView.Root
          ref={ref}
          aria-label={data.label}
          className={cx(styles.root, css(cssProps), className)}
          defaultExpandedValue={defaultExpandedValue}
          selectedValue={defaultExpandedValue}
          focusedValue={currentNode}
          {...rootProps}
        >
          <ArkTreeView.Tree className={styles.tree}>
            {sortedByVisibility.map(renderChild)}
          </ArkTreeView.Tree>
        </ArkTreeView.Root>

        {createPortal(
          <DragOverlay>
            <DragOverlayItem />
          </DragOverlay>,
          document.body,
        )}
      </DndContext>
    );
  },
);

LibraryPageTree.displayName = "DatagraphNodeTree";

type BranchProps = {
  child: NodeWithChildren;
  styles: any;
  isHighlighted: boolean;
  isRoot: boolean;
};

function DroppableBranch({
  child,
  className,
  render,
}: {
  child: NodeWithChildren;
  className: string;
  render: (props: { isOver: boolean; over: Over | null }) => JSX.Element;
}) {
  const { setNodeRef, isOver, over } = useDroppable({
    id: child.id,
  });

  return (
    <ArkTreeView.Branch
      key={child.id}
      ref={setNodeRef}
      value={child.slug}
      className={className}
    >
      {render({ isOver, over })}
    </ArkTreeView.Branch>
  );
}

const DragOverlayItem = forwardRef<any, any>(({ children, ...props }, ref) => {
  return (
    <Box
      {...props}
      ref={ref}
      borderWidth="thin"
      borderStyle="dashed"
      borderColor="gray.a6"
      height="8"
      borderRadius="md"
      backgroundColor="gray.a2"
    >
      {children}
    </Box>
  );
});
DragOverlayItem.displayName = "DragOverlayItem";

function TreeBranch({
  styles,
  child,
  isHighlighted,
  isRoot,
  isHovered,
}: BranchProps) {
  // NOTE: We need some state here to track open/close of the menu because CSS
  // isn't quite enough to track this nicely. The reason for this is that when
  // the mouse moves away from the branch control, the container that holds the
  // menu trigger moves to display: none; and the menu closes unexpectedly.
  const [menuOpen, setOpen] = useState(false);

  const isPublished = child.visibility === Visibility.published;
  const visibilityLabel = child.visibility;

  const label = isPublished ? child.name : `${child.name}`;

  const branchColourPalette = visibilityColour(child.visibility);

  const visibilityStyles = isRoot
    ? "" // Don't show the visibility state styles for root items, is cluttered.
    : css({
        paddingX: "0.5",
        borderRadius: "sm",
        colorPalette: branchColourPalette,
        borderWidth:
          child.visibility === Visibility.published ? "none" : "thin",
        borderColor: "colorPalette.8",
        borderStyle:
          child.visibility === Visibility.published ? "solid" : "dashed",
      });

  const highlightStyles = css({
    background: isHighlighted ? "gray.a2" : undefined,
  });

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: child.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    cursor: isDragging ? "grabbing" : undefined,
  };

  function handleClick(e: MouseEvent<HTMLAnchorElement>) {
    e.preventDefault();
    // if (isDragging) {
    // }
  }

  return (
    <ArkTreeView.BranchControl
      className={cx("group", styles.branchControl, highlightStyles)}
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
    >
      <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
        {child.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
      </ArkTreeView.BranchIndicator>

      <ArkTreeView.BranchText asChild className={cx(styles.branchText)}>
        <Link href={`/l/${child.slug}`} onClick={handleClick}>
          <span className={visibilityStyles}>{label}</span>{" "}
          {isHovered ? "👀" : ""}
        </Link>
      </ArkTreeView.BranchText>

      <HStack
        display={{
          base: menuOpen ? "flex" : "none",
          _groupHover: "flex",
          _active: "flex",
        }}
        gap="1"
        minW="min"
        flexShrink="0"
        onClick={() => setOpen(true)}
      >
        <CreatePageAction variant="ghost" hideLabel parentSlug={child.slug} />
        <LibraryPageMenu
          variant="ghost"
          onClose={() => setOpen(false)}
          node={child}
        />
      </HStack>
    </ArkTreeView.BranchControl>
  );
}

function TreeItem({ styles, child }: BranchProps) {
  return (
    <ArkTreeView.Item value={child.slug} className={cx(styles.item)}>
      <ArkTreeView.ItemText className={styles.itemText}>
        <Link href={`/l/${child.slug}`}>{child.name}</Link>
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
