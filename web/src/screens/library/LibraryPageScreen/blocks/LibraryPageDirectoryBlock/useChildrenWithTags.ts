import { flow, reduce, uniqBy } from "lodash/fp";

import { useNodeListChildren } from "@/api/openapi-client/nodes";
import {
  Identifier,
  NodeListResult,
  NodeWithChildren,
  TagNameList,
  TagReferenceList,
} from "@/api/openapi-schema";

// NOTE: This hook exists to pull the child list twice. Why twice? Because in
// the table block, we allow filtering by tag but the list of tags needs to be
// built from the full list of children, not the filtered subset in order to
// actually show all the tags shared by all the children.
export function useChildrenWithTags(
  nodeID: Identifier,
  initialChildren?: NodeListResult,
  childrenSort?: string,
  tagFilters?: TagNameList,
  searchQuery?: string,
) {
  const { data: filteredChildren, error: filteredError } = useNodeListChildren(
    nodeID,
    {
      children_sort: childrenSort,
      tags: tagFilters,
      q: searchQuery,
    },
    {
      swr: {
        // TODO: Perform this filtered query on the SSR too.
        fallbackData: initialChildren,
      },
    },
  );

  const { data: allChildren, error: allError } = useNodeListChildren(
    nodeID,
    undefined,
    {
      swr: {
        // TODO: Remove this when SSR has filtering.
        fallbackData: initialChildren,
      },
    },
  );

  if (!filteredChildren || !allChildren) {
    return {
      ready: false as const,
      error: filteredError || allError,
    };
  }

  const tags = getAllTags(allChildren.nodes);

  return {
    ready: true as const,
    data: filteredChildren,
    tags,
    hasChildren: allChildren.nodes.length > 0,
  };
}

const allTags = reduce<NodeWithChildren, TagReferenceList>(
  (acc, child): TagReferenceList => {
    if (child.tags) {
      return [...acc, ...child.tags];
    }
    return acc;
  },
  [],
);

const uniqueTags = uniqBy("name");

const getAllTags = flow(allTags, uniqueTags);
