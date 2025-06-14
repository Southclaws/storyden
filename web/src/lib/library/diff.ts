import { dequal } from "dequal";

import { NodeMutableProps, NodeWithChildren } from "@/api/openapi-schema";

function projectNodeToMutableProps(node: NodeWithChildren): NodeMutableProps {
  return {
    name: node.name,
    slug: node.slug,
    content: node.content,
    tags: node.tags.map((t) => t.name),
    primary_image_asset_id: node.primary_image?.id,
    url: node.link?.url,
    properties: node.properties,
    hide_child_tree: node.children.length === 0,
    meta: node.meta,
  };
}

export const deriveMutationFromDifference = (
  current: NodeWithChildren,
  updated: NodeWithChildren,
) => {
  const mutation: NodeMutableProps = {};

  const draft = projectNodeToMutableProps(current);
  const changes = projectNodeToMutableProps(updated);

  (Object.keys(changes) as (keyof NodeMutableProps)[]).forEach((key) => {
    const draftValue = draft[key];
    const updatedValue = changes[key];
    if (!updatedValue) {
      return;
    }

    const changed = !dequal(draftValue, updatedValue);
    if (!changed) {
      return;
    }

    Object.assign(mutation, { [key]: updatedValue });
  });

  return mutation;
};
