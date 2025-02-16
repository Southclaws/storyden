import { uniqueId } from "lodash/fp";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getNodeGetKey,
  nodeUpdateChildrenPropertySchema,
  nodeUpdatePropertySchema,
} from "@/api/openapi-client/nodes";
import {
  Identifier,
  Node,
  NodeGetOKResponse,
  NodeWithChildren,
  PropertyList,
  PropertySchema,
  PropertySchemaList,
  PropertySchemaMutableProps,
} from "@/api/openapi-schema";

export function usePropertyMutation(node: NodeWithChildren) {
  const { mutate } = useSWRConfig();

  const nodeKey = node && getNodeGetKey(node.slug);
  const nodeKeyFn =
    node &&
    ((key: Arguments) => {
      return Array.isArray(key) && key[0].startsWith(nodeKey);
    });

  const addField = async (newProperty: PropertySchemaMutableProps) => {
    await handle(async () => {
      const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
        if (!data) return;

        const newPropertyList = [
          ...node.properties,
          {
            fid: uniqueId("optimistic_property_mutation_"),
            ...newProperty,
          },
        ] satisfies PropertyList;

        const updated = {
          ...data,
          properties: newPropertyList,
        } satisfies NodeWithChildren;

        return updated;
      };

      mutate(nodeKeyFn, nodeMutator);

      await nodeUpdatePropertySchema(node.slug, [
        ...node.properties,
        newProperty,
      ]);
    });
  };

  const removeField = async (fid: Identifier) => {
    await handle(async () => {
      const newPropertyList = node.properties.filter((p) => p.fid !== fid);

      const nodeMutator: MutatorCallback<NodeGetOKResponse> = (data) => {
        if (!data) return;

        const updated = {
          ...data,
          properties: newPropertyList,
        } satisfies NodeWithChildren;

        return updated;
      };

      mutate(nodeKeyFn, nodeMutator);

      await nodeUpdatePropertySchema(node.slug, newPropertyList);
    });
  };

  return {
    addField,
    removeField,
  };
}
