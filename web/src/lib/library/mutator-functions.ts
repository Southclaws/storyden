import { mergeWith } from "lodash";
import { MutatorCallback } from "swr";

import {
  Node,
  NodeGetOKResponse,
  NodeListOKResponse,
  NodeWithChildren,
} from "@/api/openapi-schema";

export const nodeListMutator = (
  updated: NodeWithChildren,
): MutatorCallback<NodeListOKResponse> => {
  return (data) => {
    if (!data) return;

    const nodes = data.nodes.map((node) => {
      if (node.slug === updated.slug) {
        return {
          ...node,
          ...updated,
          // Ensure the children are not included if the node is supposed to
          // hide them. The reason this needs to be done explicitly is because
          // these mutator functions are often called with a node state that was
          // pulled from GET /nodes/{id} API which always includes children, not
          // like the tree listing API which will adhere to child-hiding rules.
          children: updated.hide_child_tree ? [] : updated.children,
        };
      }

      return node;
    });

    return {
      ...data,
      nodes,
    };
  };
};

export const nodeMutator = (
  updated: NodeWithChildren,
): MutatorCallback<NodeGetOKResponse> => {
  return (previous) => {
    if (!previous) return;

    const node = mergeWith(previous, updated, (target, source) => {
      // If the target is an array, we want to replace it with the source value.
      if (Array.isArray(target)) {
        return source;
      }

      // If source is undefined, we want to keep the existing target value.
      if (source === undefined) {
        return target;
      }

      return source;
    });

    return node;
  };
};
