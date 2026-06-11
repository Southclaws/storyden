import { pick } from "lodash";

import {
  NodeVersion,
  NodeVersionMutableProps,
  NodeWithChildren,
  PropertyMutationList,
  PropertyType,
} from "@/api/openapi-schema";
import { MutationSet } from "@/lib/library/diff";
import { WithMetadata, hydrateNode } from "@/lib/library/metadata";

export function buildNodeVersionMutation(
  mutation: MutationSet,
): NodeVersionMutableProps {
  return pick(mutation.nodeMutation, [
    "content",
    "description",
    "name",
    "properties",
    "slug",
  ]);
}

export function overlayNodeVersion(
  node: WithMetadata<NodeWithChildren>,
  version: NodeVersion,
): WithMetadata<NodeWithChildren> {
  const draft = structuredClone(node);

  draft.name = version.name;
  draft.slug = version.slug;
  draft.properties = normaliseVersionProperties(version.properties);

  if (version.description !== undefined) {
    draft.description = version.description ?? "";
  }

  if (version.content !== undefined) {
    draft.content = version.content ?? undefined;
  }

  return hydrateNode(draft);
}

function normaliseVersionProperties(properties: PropertyMutationList) {
  return properties.map((property, index) => ({
    fid: property.fid ?? `new_field_${index}`,
    name: property.name,
    sort: property.sort ?? `${index}`,
    type: property.type ?? PropertyType.text,
    value: property.value,
  }));
}
