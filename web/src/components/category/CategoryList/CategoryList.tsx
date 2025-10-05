"use client";

import {
  TreeView as ArkTreeView,
  createTreeCollection,
} from "@ark-ui/react/tree-view";
import { useDraggable, useDroppable } from "@dnd-kit/core";
import { SortableContext, rectSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { CSSProperties, useEffect, useMemo, useRef, useState } from "react";
import { KeyedMutator } from "swr";

import { handle } from "@/api/client";
import {
  categoryUpdatePosition,
  useCategoryList as useGetCategoryList,
} from "@/api/openapi-client/categories";
import {
  Category,
  CategoryListOKResponse,
  CategoryUpdatePositionBody,
  Identifier,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { CategoryMenu } from "@/components/category/CategoryMenu/CategoryMenu";
import { Anchor } from "@/components/site/Anchor";
import { DiscussionRoute } from "@/components/site/Navigation/Anchors/Discussion";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { Unready } from "@/components/site/Unready";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { ChevronRightIcon } from "@/components/ui/icons/Chevron";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { useCategoryEvent } from "@/lib/category/events";
import {
  CategoryTree,
  buildCategoryTree,
  isDescendant,
} from "@/lib/category/tree";
import {
  DragItemCategory,
  DragItemCategoryDivider,
  DragItemData,
} from "@/lib/dragdrop/provider";
import { css, cx } from "@/styled-system/css";
import { HStack, LStack } from "@/styled-system/jsx";
import { treeView } from "@/styled-system/recipes";
import { hasPermission } from "@/utils/permissions";

import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";

export type Props = {
  initialCategoryList?: CategoryListOKResponse;
  currentCategorySlug?: string;
};

type PositionInList = "top" | "in" | "bottom" | "only";

export function CategoryList({
  initialCategoryList,
  currentCategorySlug,
}: Props) {
  const { data, error, mutate } = useGetCategoryList({
    swr: { fallbackData: initialCategoryList },
  });

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <CategoryListTree
      categories={data.categories}
      currentCategorySlug={currentCategorySlug}
      mutate={mutate}
    />
  );
}

export function CategoryListTree({
  categories,
  currentCategorySlug,
  mutate,
}: {
  categories: Category[];
  currentCategorySlug?: string;
  mutate: KeyedMutator<CategoryListOKResponse>;
}) {
  const session = useSession();

  const canManageCategories = hasPermission(session, "MANAGE_CATEGORIES");

  const tree = buildCategoryTree(categories);

  const collection = createTreeCollection<CategoryTree>({
    nodeToValue: (cat) => cat?.id ?? "",
    nodeToString: (cat) => cat?.slug ?? "",
    nodeToChildren: (cat) => cat?.children ?? [],
    rootNode: { children: tree } as CategoryTree,
  });

  const rootNodes = collection.rootNode.children;

  const defaultExpanded = useMemo(() => {
    return rootNodes.map((node) => node.id);
  }, [rootNodes]);

  const [expandedValue, setExpandedValue] = useState<string[]>(defaultExpanded);

  const handleExpandedChange = (details: { expandedValue: string[] }) => {
    setExpandedValue(details.expandedValue);
  };

  const handleExpandNode = (id: string) => {
    setExpandedValue((prev) => (prev.includes(id) ? prev : [...prev, id]));
  };

  useCategoryEvent(
    "category:reorder-category",
    async ({ categorySlug, direction, newParent, targetCategory }) => {
      await handle(async () => {
        const params: CategoryUpdatePositionBody = (() => {
          switch (direction) {
            case "above":
              return {
                before: targetCategory,
                parent: newParent,
              };

            case "below":
              return {
                after: targetCategory,
                parent: newParent,
              };

            case "inside":
              return {
                parent: targetCategory,
              };
          }
        })();

        const response = await categoryUpdatePosition(categorySlug, params);
        await mutate(response, { revalidate: false });
      });
    },
  );

  if (!tree) {
    return <Unready />;
  }

  const styles = treeView();

  return (
    <LStack gap="0">
      <NavigationHeader
        href={DiscussionRoute}
        controls={canManageCategories && <CategoryCreateTrigger hideLabel />}
      >
        <HStack gap="1">
          <DiscussionIcon />
          Discussion
        </HStack>
      </NavigationHeader>

      <ArkTreeView.Root
        className={styles.root}
        collection={collection}
        expandedValue={expandedValue}
        onExpandedChange={handleExpandedChange}
      >
        <ArkTreeView.Tree className={styles.tree}>
          <SortableContext
            items={rootNodes.map((child) => child.id)}
            strategy={rectSortingStrategy}
          >
            {rootNodes.map((cat, index) => (
              <CategoryTreeNode
                key={cat.id}
                fullTree={tree}
                currentCategorySlug={currentCategorySlug}
                parentID={null}
                category={cat}
                styles={styles}
                isRoot={true}
                indexPath={[]}
                positionInList={getPositionInList(rootNodes.length, index)}
                handleExpandNode={handleExpandNode}
                canManageCategories={canManageCategories}
              />
            ))}
          </SortableContext>
        </ArkTreeView.Tree>
      </ArkTreeView.Root>
    </LStack>
  );
}

type TreeNodeProps = {
  fullTree: CategoryTree[];
  currentCategorySlug: string | undefined;
  parentID: Identifier | null;
  category: CategoryTree;
  styles: any;
  isRoot: boolean;
  indexPath: number[];
  positionInList: PositionInList;
  handleExpandNode: (id: string) => void;
  canManageCategories: boolean;
};

function CategoryTreeNode({
  fullTree,
  currentCategorySlug,
  parentID,
  category,
  styles,
  isRoot,
  indexPath,
  positionInList,
  handleExpandNode,
  canManageCategories,
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
    id: category.id,
    disabled: !canManageCategories,
    data: {
      type: "category",
      categoryID: category.id,
      category: category,
      hasChildren: category.children.length > 0,
    } satisfies DragItemCategory,
  });

  const dragged = active?.data.current as DragItemCategory | undefined;
  const draggedID = dragged?.categoryID;

  const isDescendantOfDragged = draggedID
    ? isDescendant(fullTree, draggedID, category.id)
    : false;

  const { setNodeRef: setDroppableNodeRef } = useDroppable({
    id: category.id,
    disabled:
      !canManageCategories ||
      isInvalidDropTarget(draggedID ?? null, category.id, fullTree),
    data: {
      type: "category",
      categoryID: category.id,
      category: category,
      hasChildren: category.children.length > 0,
    } satisfies DragItemCategory,
  });

  const overItem = over?.data.current as DragItemData | undefined;

  const isDraggingOver =
    !isDescendantOfDragged &&
    overItem?.type === "category" &&
    dragged?.type === "category" &&
    overItem?.categoryID === category.id &&
    !isDragging;

  // handle drag-over to expand
  const expandTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  useEffect(() => {
    if (!canManageCategories) return;

    if (isDraggingOver) {
      expandTimeout.current = setTimeout(() => {
        handleExpandNode(category.id);
      }, 600);
    } else if (expandTimeout.current) {
      clearTimeout(expandTimeout.current);
      expandTimeout.current = null;
    }

    return () => {
      if (expandTimeout.current) {
        clearTimeout(expandTimeout.current);
        expandTimeout.current = null;
      }
    };
  }, [canManageCategories, handleExpandNode, isDraggingOver, category.id]);

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
    transform: CSS.Transform.toString(transform),
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
    overItem?.type === "category-divider" &&
    overItem?.siblingCategoryID === category.id &&
    overItem?.direction === "above";
  const showBottomDivider =
    overItem?.type === "category-divider" &&
    overItem?.siblingCategoryID === category.id &&
    overItem?.direction === "below";

  const isHighlighted = category.slug === currentCategorySlug;

  const draggingOverStyles = isDraggingOver
    ? css({
        borderRadius: "md",
        colorPalette: "accent",
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
    <ArkTreeView.NodeProvider
      key={category.id}
      node={category}
      indexPath={indexPath}
    >
      <DropIndicator
        categoryID={category.id}
        parentCategoryID={category.parent}
        active={showTopDivider}
        direction="above"
        positionInList={positionInList}
        canManageCategories={canManageCategories}
      />

      <ArkTreeView.Branch
        ref={setDroppableNodeRef}
        className={cx(
          styles.branch,
          // branchIndentClass,
          // css({
          //   // "&[data-drag-over=true]": {
          //   //   outline: "1px dashed var(--colors-border.emphasized)",
          //   // },
          // }),
        )}
        data-drag-over={isDraggingOver}
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
            {category.children.length > 0 ? (
              <ChevronRightIcon />
            ) : (
              <BulletIcon />
            )}
          </ArkTreeView.BranchIndicator>

          <ArkTreeView.BranchText asChild className={styles.branchText}>
            <Anchor
              href={`/d/${category.slug}`}
              className={css({
                display: "flex",
                alignItems: "center",
                gap: "2",
                flex: "1",
                _hover: { textDecoration: "none" },
              })}
            >
              {category.name}
            </Anchor>
          </ArkTreeView.BranchText>

          {canManageCategories && (
            <HStack
              gap="1"
              opacity={{
                base: "0",
                _groupHover: "full",
              }}
            >
              <CategoryMenu category={category} />
            </HStack>
          )}
        </ArkTreeView.BranchControl>

        <ArkTreeView.BranchContent className={styles.branchContent}>
          <SortableContext
            items={category.children.map((child) => child.id)}
            strategy={rectSortingStrategy}
          >
            {category.children.map((child, childIndex) => (
              <CategoryTreeNode
                key={child.id}
                fullTree={fullTree}
                currentCategorySlug={currentCategorySlug}
                parentID={category.parent ?? null}
                category={child}
                indexPath={[...indexPath, childIndex]}
                isRoot={false}
                styles={styles}
                positionInList={getPositionInList(
                  category.children.length,
                  childIndex,
                )}
                handleExpandNode={handleExpandNode}
                canManageCategories={canManageCategories}
              />
            ))}
          </SortableContext>
        </ArkTreeView.BranchContent>
      </ArkTreeView.Branch>

      <DropIndicator
        categoryID={category.id}
        parentCategoryID={category.parent}
        direction="below"
        active={showBottomDivider}
        positionInList={positionInList}
        canManageCategories={canManageCategories}
      />
    </ArkTreeView.NodeProvider>
  );
}

type DropIndicatorProps = {
  categoryID: string;
  parentCategoryID: string | undefined;
  direction: "above" | "below";
  active: boolean;
  positionInList: PositionInList;
  canManageCategories: boolean;
};

function DropIndicator({
  categoryID,
  parentCategoryID,
  direction,
  active,
  positionInList,
  canManageCategories,
}: DropIndicatorProps) {
  const { setNodeRef } = useDroppable({
    id: `${categoryID}_${direction}`,
    disabled: false,
    data: {
      type: "category-divider",
      direction,
      siblingCategoryID: categoryID,
      parentID: parentCategoryID ?? null,
    } satisfies DragItemCategoryDivider,
  });

  if (!canManageCategories) {
    return null;
  }

  const shouldRender =
    direction === "above"
      ? positionInList === "top" ||
        positionInList === "in" ||
        positionInList === "only"
      : positionInList === "bottom" ||
        positionInList === "in" ||
        positionInList === "only";

  if (!shouldRender) {
    return null;
  }

  return (
    <div
      data-divider-category-id={categoryID}
      data-divider-parent-category-id={parentCategoryID}
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

function isInvalidDropTarget(
  activeId: string | null,
  targetId: string,
  tree: CategoryTree[],
) {
  if (!activeId) {
    return false;
  }

  if (activeId === targetId) {
    return true;
  }

  return isDescendant(tree, activeId, targetId);
}

function getPositionInList(length: number, index: number): PositionInList {
  if (length === 1) {
    return "only";
  }

  if (index === 0) {
    return "top";
  }

  if (index === length - 1) {
    return "bottom";
  }

  return "in";
}
