import { omit } from "lodash";

import {
  Identifier,
  NodeMutableProps,
  NodeWithChildren,
  PropertySchemaList,
  PropertySchemaMutableProps,
} from "@/api/openapi-schema";
import { deepEqual } from "@/utils/equality";

import { isValidLinkLike, normalizeLink } from "../link/validation";
import { isSlugReady } from "../mark/mark";

type NodeMutablePropsWithChildren = NodeMutableProps &
  Pick<NodeWithChildren, "children">;

function projectNodeToMutableProps(
  node: NodeWithChildren,
): NodeMutablePropsWithChildren {
  return {
    name: node.name,
    slug: node.slug,
    content: node.content,
    tags: node.tags.map((t) => t.name),
    primary_image_asset_id: node.primary_image?.id,
    url: normalizeLink(node.link?.url),
    description: node.description,
    properties: node.properties.map((p) => {
      const fid = p.fid.startsWith("new_field") ? undefined : p.fid;

      return {
        fid: fid,
        name: p.name,
        value: p.value,
        type: p.type,
      };
    }),
    hide_child_tree: node.hide_child_tree,
    meta: node.meta,
    children: node.children,
  };
}

export type MutationSet = {
  clean: boolean;
  nodeMutation: NodeMutableProps;
  childPropertySchemaMutation?: PropertySchemaMutableProps[];
  childMutation: Record<Identifier, NodeMutableProps[]>;
};

export const deriveMutationFromDifference = (
  current: NodeWithChildren,
  updated: NodeWithChildren,
): MutationSet => {
  const mutation: NodeMutableProps = {};
  const childMutation: Record<Identifier, NodeMutableProps[]> = {};

  const draft = projectNodeToMutableProps(current);
  const changes = projectNodeToMutableProps(updated);

  const keys = Object.keys(changes) as (keyof NodeMutablePropsWithChildren)[];

  console.debug(
    "deriveMutationFromDifference",
    current.id,
    "keys",
    keys,
    "changes",
    changes,
  );

  keys.forEach((key) => {
    const draftValue = draft[key];
    const updatedValue = changes[key];
    if (updatedValue === null) {
      console.debug(`Skipping mutation for '${key}' because it is null`);
      return;
    }

    const changed = !deepEqual(draftValue, updatedValue);
    if (!changed) {
      console.debug(
        `Skipping mutation for '${key}' because it has not changed`,
        {
          old: draftValue,
          new: updatedValue,
        },
      );
      return;
    }

    // Field specific transformations and skipping logic.

    switch (key) {
      case "slug": {
        const slugValue = updatedValue as string | undefined;
        if (slugValue === undefined) {
          return;
        }
        if (!isSlugReady(slugValue)) {
          // Slugs must be valid to be added to patch, see mark.ts for details.
          console.debug("Skipping mutation for 'slug' because it is not ready");
          return;
        }
        break;
      }
      case "url": {
        if (updated.link === undefined) {
          Object.assign(mutation, { url: null });
          return;
        }

        const rawURL = updated.link.url;
        if (!isValidLinkLike(rawURL)) {
          console.debug(
            "Skipping mutation for 'url' because it is not valid",
            rawURL,
          );
          return;
        }

        const normalizedURL = normalizeLink(rawURL);
        if (normalizedURL === undefined) {
          console.debug(
            "Skipping mutation for 'url' because it could not be normalized",
            rawURL,
          );
          return;
        }

        Object.assign(mutation, { url: normalizedURL });
        return;
      }
      case "primary_image_asset_id": {
        const primaryImageAssetId = updatedValue as Identifier | undefined;
        if (primaryImageAssetId === undefined) {
          Object.assign(mutation, { primary_image_asset_id: null });
          return;
        }
        break;
      }
      case "children": {
        const updatedChildren = updatedValue as NodeWithChildren["children"];
        if (updatedChildren.length === 0) {
          console.debug("Skipping mutation for 'children' because it is empty");
          return;
        }

        // If children have changed we need to create a mutation for each child.
        updatedChildren.forEach((child) => {
          const childDraft = draft.children.find(
            (c) => c.id === child.id,
          ) as NodeWithChildren;

          if (!childDraft) {
            console.warn("Child draft not found for", child.id);
            return;
          }

          // Safe recursion: children do not contain a full tree of descendants.
          const childChanges = deriveMutationFromDifference(childDraft, child);
          console.debug("child deriveMutationFromDifference:", childChanges);
          if (childChanges.clean) {
            return;
          }

          console.debug(
            "Adding child mutation for",
            child.id,
            childChanges.nodeMutation,
          );

          (childMutation[child.id] ??= []).push(childChanges.nodeMutation);
        });

        return;
      }
    }

    Object.assign(mutation, { [key]: updatedValue });
  });

  // Diff the child property schema
  const childPropertySchema = diffPropertySchemas(
    current.child_property_schema,
    updated.child_property_schema,
  );

  const nodeMutations = Object.keys(mutation).length;

  // Determine if this mutation even does anything.
  const clean =
    nodeMutations === 0 &&
    !childPropertySchema &&
    !Object.keys(childMutation).length;

  return {
    clean,
    nodeMutation: mutation,
    childPropertySchemaMutation: childPropertySchema,
    childMutation: childMutation,
  };
};

function diffPropertySchemas(
  a: PropertySchemaList,
  b: PropertySchemaList,
): PropertySchemaMutableProps[] | undefined {
  if (deepEqual(a, b)) {
    return undefined;
  }

  return b.map((p) => {
    if (p.fid.startsWith("new_field")) {
      return omit(p, "fid") as PropertySchemaMutableProps;
    }

    return p;
  });
}
