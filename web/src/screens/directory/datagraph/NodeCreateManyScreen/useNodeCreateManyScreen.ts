import { last } from "lodash";

import { nodeCreate, useNodeGet } from "src/api/openapi-client/nodes";
import { Link, LinkReference, NodeWithChildren } from "src/api/openapi-schema";
import { DatagraphNodeWithRelations } from "src/components/directory/datagraph/DatagraphNode";

import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  node?: NodeWithChildren;
};

export function useNodeCreateManyScreen(props: Props) {
  const directoryPath = useDirectoryPath();
  const { data } = useNodeGet(props.node?.slug ?? (null as any as string), {
    swr: {
      fallbackData: props.node,
    },
  });

  async function handleCreate(
    link: LinkReference,
  ): Promise<DatagraphNodeWithRelations> {
    const parentSlug = last(directoryPath as string[]);
    const created = await nodeCreate({
      name: link.title || link.url,
      slug: link.slug,
      url: link.url,
      content: link.description,
      parent: parentSlug,
      visibility: "draft",
    });

    return {
      ...created,
      children: [],
      recomentations: [],
    };
  }

  return {
    parent: data,
    handlers: {
      handleCreate,
    },
  };
}
