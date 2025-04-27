import {
  TreeView as ArkTreeView,
  createTreeCollection,
} from "@ark-ui/react/tree-view";
import Link from "next/link";
import { Fragment, JSX } from "react";

import { NodeWithChildren, Visibility } from "@/api/openapi-schema";
import { CreatePageAction } from "@/components/library/CreatePage";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { ReviewIcon, UnlistedIcon } from "@/components/ui/icons/Visibility";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css, cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";

import { LibraryPageMenu } from "../LibraryPageMenu/LibraryPageMenu";

export interface LibraryPageTreeProps {
  nodes: NodeWithChildren[];
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
    nodeToValue: (n: NodeWithChildren) => {
      return n.id;
    },
    nodeToString: (n: NodeWithChildren) => {
      return n.slug;
    },
    nodeToChildren: (n: NodeWithChildren) => {
      return n.children;
    },
    rootNode: {
      children: sortedByVisibility,
    } as NodeWithChildren,
  });

  const rootNodes = collection.rootNode.children;

  return (
    <ArkTreeView.Root
      className={styles.root}
      collection={collection}
      defaultExpandedValue={defaultExpandedValue}
    >
      <ArkTreeView.Tree className={cx(styles.tree)}>
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
                currentNode={currentNode}
                node={node}
                indexPath={[]}
                isRoot={true}
                styles={styles}
              />
            </Fragment>
          );
        })}
      </ArkTreeView.Tree>
    </ArkTreeView.Root>
  );
};

type TreeNodeProps = {
  currentNode: string | undefined;
  node: NodeWithChildren;
  styles: any;
  isRoot: boolean;
  indexPath: number[];
};

function TreeNode({
  currentNode,
  styles,
  node,
  isRoot,
  indexPath,
}: TreeNodeProps) {
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
        borderColor: "colorPalette.8",
        borderStyle:
          node.visibility === Visibility.published ? "solid" : "dashed",
      });

  const highlightStyles = css({
    background: isHighlighted ? "gray.a2" : undefined,
  });

  return (
    <ArkTreeView.NodeProvider
      key={node.id}
      node={node}
      indexPath={indexPath}
      data-visibility={node.visibility}
    >
      <ArkTreeView.Branch className={cx(styles.branch)}>
        <ArkTreeView.BranchControl
          className={cx("group", styles.branchControl, highlightStyles)}
        >
          <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
            {node.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
          </ArkTreeView.BranchIndicator>

          <ArkTreeView.BranchText
            asChild
            className={cx(styles.branchText, visibilityStyles)}
          >
            <Link href={`/l/${node.slug}`}>{label}</Link>
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
            <CreatePageAction
              variant="ghost"
              hideLabel
              parentSlug={node.slug}
            />
            <LibraryPageMenu variant="ghost" node={node} />
          </HStack>
        </ArkTreeView.BranchControl>

        <ArkTreeView.BranchContent className={styles.branchContent}>
          {node.children.map((child, index) => {
            return (
              <TreeNode
                key={child.id}
                currentNode={currentNode}
                node={child}
                indexPath={[...indexPath, index]}
                isRoot={false}
                styles={styles}
              />
            );
          })}
        </ArkTreeView.BranchContent>
      </ArkTreeView.Branch>
    </ArkTreeView.NodeProvider>
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
