import { dequal } from "dequal";
import { omit } from "lodash";

import {
  NodeMutableProps,
  NodeWithChildren,
  PropertySchemaList,
  PropertySchemaMutableProps,
} from "@/api/openapi-schema";

import { isSlugReady } from "../mark/mark";

function projectNodeToMutableProps(node: NodeWithChildren): NodeMutableProps {
  return {
    name: node.name,
    slug: node.slug,
    content: node.content,
    tags: node.tags.map((t) => t.name),
    primary_image_asset_id: node.primary_image?.id,
    url: node.link?.url,
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
  };
}

export type MutationSet = {
  clean: boolean;
  nodeMutation: NodeMutableProps;
  childPropertySchemaMutation?: PropertySchemaMutableProps[];
};

export const deriveMutationFromDifference = (
  current: NodeWithChildren,
  updated: NodeWithChildren,
): MutationSet => {
  const mutation: NodeMutableProps = {};

  const draft = projectNodeToMutableProps(current);
  const changes = projectNodeToMutableProps(updated);

  (Object.keys(changes) as (keyof NodeMutableProps)[]).forEach((key) => {
    const draftValue = draft[key];
    const updatedValue = changes[key];
    if (updatedValue === undefined || updatedValue === null) {
      return;
    }

    const changed = !dequal(draftValue, updatedValue);
    if (!changed) {
      return;
    }

    // Field specific transformations and skipping logic.

    switch (key) {
      case "slug":
        if (!isSlugReady(updatedValue as string)) {
          // Slugs must be valid to be added to patch, see mark.ts for details.
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
  const clean = nodeMutations === 0 && !childPropertySchema;

  return {
    clean,
    nodeMutation: mutation,
    childPropertySchemaMutation: childPropertySchema,
  };
};

function diffPropertySchemas(
  a: PropertySchemaList,
  b: PropertySchemaList,
): PropertySchemaMutableProps[] | undefined {
  if (dequal(a, b)) {
    return undefined;
  }

  return b.map((p) => {
    if (p.fid.startsWith("new_field")) {
      return omit(p, "fid") as PropertySchemaMutableProps;
    }

    return p;
  });
}
