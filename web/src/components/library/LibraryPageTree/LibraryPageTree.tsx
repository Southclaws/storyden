import type { Assign } from "@ark-ui/react";
import {
  TreeView as ArkTreeView,
  type TreeViewRootProps,
} from "@ark-ui/react/tree-view";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { keyBy } from "lodash";
import Link from "next/link";
import { JSX, forwardRef } from "react";

import { NodeWithChildren, Visibility } from "@/api/openapi-schema";
import { CreatePageAction } from "@/components/library/CreatePage";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { ReviewIcon, UnlistedIcon } from "@/components/ui/icons/Visibility";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx } from "@/styled-system/css";
import { HStack, splitCssProps } from "@/styled-system/jsx";
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

const visibilityIcons: Record<Visibility, JSX.Element> = {
  [Visibility.published]: <></>,
  [Visibility.review]: <ReviewIcon />,
  [Visibility.draft]: <DraftIcon />,
  [Visibility.unlisted]: <UnlistedIcon />,
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
      const dividerIcon = visibilityIcons[child.visibility];

      // We only show dividers on the root list, as this is the only list that's
      // sorted by visibility.
      const showDivider = isRoot && !sameVisibilityAsPrevious;

      const isHighlighted = child.slug === currentNode;

      return (
        <ArkTreeView.Branch
          key={child.id}
          value={child.slug}
          className={styles.branch}
        >
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

          <TreeBranch
            styles={styles}
            child={child}
            isHighlighted={isHighlighted}
            isRoot={isRoot}
          />

          <ArkTreeView.BranchContent className={styles.branchContent}>
            {child.children?.map((child, i) =>
              child.children ? (
                renderChild(child, i)
              ) : (
                <TreeItem
                  key={child.id}
                  styles={styles}
                  child={child}
                  isHighlighted={isHighlighted}
                />
              ),
            )}
          </ArkTreeView.BranchContent>
        </ArkTreeView.Branch>
      );
    };

    const sortedByVisibility = data.children.sort((a, b) => {
      return visibilitySortKey[a.visibility] - visibilitySortKey[b.visibility];
    });

    return (
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
    );
  },
);

LibraryPageTree.displayName = "DatagraphNodeTree";

type BranchProps = {
  child: NodeWithChildren;
  styles: any;
  isHighlighted: boolean;
  isRoot?: boolean;
};

function TreeBranch({ styles, child, isHighlighted, isRoot }: BranchProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: child.id,
  });

  const isPublished = child.visibility === Visibility.published;
  const label = isPublished ? child.name : `${child.name}`;
  const branchColourPalette = visibilityColour(child.visibility);

  const visibilityStyles = isRoot
    ? ""
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
    colorPalette: branchColourPalette,
    backgroundColor: isHighlighted ? "colorPalette.2" : undefined,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : undefined,
  };

  return (
    <ArkTreeView.BranchControl
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={cx(
        styles.branchControl,
        css({
          cursor: "grab",
          _active: {
            cursor: "grabbing",
          },
        }),
        visibilityStyles,
        highlightStyles,
      )}
    >
      <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
        <ChevronRightIcon />
      </ArkTreeView.BranchIndicator>

      <ArkTreeView.BranchText asChild className={styles.branchText}>
        <Link href={`/library/${child.slug}`}>
          {label}
        </Link>
      </ArkTreeView.BranchText>

      <LibraryPageMenu node={child} />
    </ArkTreeView.BranchControl>
  );
}

function TreeItem({ styles, child, isHighlighted }: BranchProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: child.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : undefined,
  };

  return (
    <ArkTreeView.Item
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      value={child.slug}
      className={cx(
        styles.item,
        css({
          cursor: "grab",
          _active: {
            cursor: "grabbing",
          },
        }),
      )}
    >
      <ArkTreeView.ItemText className={styles.itemText}>
        <Link href={`/library/${child.slug}`}>{child.name}</Link>
      </ArkTreeView.ItemText>
      <LibraryPageMenu node={child} />
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
