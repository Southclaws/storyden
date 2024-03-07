import { last } from "lodash";

import { clusterCreate, useClusterGet } from "src/api/openapi/clusters";
import { ClusterWithItems, Link } from "src/api/openapi/schemas";
import { DatagraphNodeWithRelations } from "src/components/directory/datagraph/DatagraphNode";

import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  cluster?: ClusterWithItems;
};

export function useClusterCreateManyScreen(props: Props) {
  const directoryPath = useDirectoryPath();
  const { data } = useClusterGet(
    props.cluster?.slug ?? (null as any as string),
    {
      swr: {
        fallbackData: props.cluster,
      },
    },
  );

  async function handleCreate(link: Link): Promise<DatagraphNodeWithRelations> {
    const parentSlug = last(directoryPath as string[]);
    const created = await clusterCreate({
      name: link.title || link.url,
      slug: link.slug,
      url: link.url,
      description: link.description || "(No description)",
      asset_ids: link.assets.map((asset) => asset.id),
      parent: parentSlug,
      visibility: "draft",
    });

    return {
      ...created,
      type: "cluster",
      clusters: [],
      items: [],
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
