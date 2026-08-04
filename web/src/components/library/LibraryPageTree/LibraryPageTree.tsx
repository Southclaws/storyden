import { Fragment, JSX, useMemo } from "react";

import { NodeWithChildren, Visibility } from "@/api/openapi-schema";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { ReviewIcon, UnlistedIcon } from "@/components/ui/icons/Visibility";
import { DragTree } from "@/components/ui/tree-view";
import { DragItemDivider, DragItemNode } from "@/lib/dragdrop/provider";
import { visibilityColour } from "@/lib/library/visibilityColours";
import { css } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

import { CreatePageAction } from "../CreatePage";
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

const getNodeId = (node: NodeWithChildren) => node.id;
const getNodeLabel = (node: NodeWithChildren) => node.name;
const getNodeHref = (node: NodeWithChildren) => `/l/${node.slug}`;
const getNodeChildren = (node: NodeWithChildren) => node.children;

export function LibraryPageTree({
  nodes,
  currentNode,
  canManageLibrary,
}: LibraryPageTreeProps) {
  const sortedNodes = useMemo(
    () =>
      [...nodes].sort(
        (a, b) =>
          visibilitySortKey[a.visibility] - visibilitySortKey[b.visibility],
      ),
    [nodes],
  );
  const expandedIds = useMemo(
    () => getExpandedNodeIds(sortedNodes, currentNode),
    [currentNode, sortedNodes],
  );

  return (
    <DragTree
      id="library-navigation"
      label="Library pages"
      nodes={sortedNodes}
      getId={getNodeId}
      getLabel={getNodeLabel}
      getHref={getNodeHref}
      getChildren={getNodeChildren}
      defaultExpandedIds={expandedIds}
      isSelected={(node) => node.slug === currentNode}
      drag={{
        enabled: canManageLibrary,
        getNodeData: ({ node, parent }) =>
          ({
            type: "node",
            node,
            parentID: parent?.id ?? null,
            context: "sidebar",
          }) satisfies DragItemNode,
        getDividerData: ({ node, parent, direction }) =>
          ({
            type: "divider",
            direction,
            siblingNode: node,
            parentID: parent?.id ?? null,
            context: "sidebar",
          }) satisfies DragItemDivider,
      }}
      getLabelClassName={({ node, isRoot }) =>
        isRoot ? undefined : getVisibilityClassName(node)
      }
      renderBefore={({ node, index, siblings, isRoot }) => {
        if (!isRoot || index === 0) {
          return null;
        }
        const previous = siblings[index - 1];
        if (!previous || previous.visibility === node.visibility) {
          return null;
        }

        return (
          <HStack mb="1">
            <NavigationHeader href="/drafts">
              <HStack>
                {visibilityIcons[node.visibility]}
                {visibilityLabels[node.visibility]}
              </HStack>
            </NavigationHeader>
          </HStack>
        );
      }}
      renderActions={({ node }) => (
        <Fragment>
          {canManageLibrary && (
            <CreatePageAction
              variant="plain"
              hideLabel
              parentSlug={node.slug}
            />
          )}
          <LibraryPageMenu variant="plain" node={node} />
        </Fragment>
      )}
    />
  );
}

function getVisibilityClassName(node: NodeWithChildren) {
  const published = node.visibility === Visibility.published;

  return css({
    borderColor: published ? "transparent" : "colorPalette.border",
    borderRadius: "sm",
    borderStyle: published ? "solid" : "dashed",
    borderWidth: published ? "none" : "thin",
    colorPalette: visibilityColour(node.visibility),
    paddingX: "0.5",
  });
}

function getExpandedNodeIds(
  nodes: NodeWithChildren[],
  currentSlug: string | undefined,
) {
  if (!currentSlug) {
    return [];
  }

  const visit = (node: NodeWithChildren): string[] | null => {
    if (node.slug === currentSlug) {
      return [node.id];
    }

    for (const child of node.children) {
      const path = visit(child);
      if (path) {
        return [node.id, ...path];
      }
    }

    return null;
  };

  for (const node of nodes) {
    const path = visit(node);
    if (path) {
      return path;
    }
  }

  return [];
}
