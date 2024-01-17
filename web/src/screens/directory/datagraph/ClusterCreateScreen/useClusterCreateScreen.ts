import { useRouter } from "next/navigation";

import { clusterCreate } from "src/api/openapi/clusters";
import {
  Account,
  ClusterInitialProps,
  ClusterWithItems,
} from "src/api/openapi/schemas";

import { joinDirectoryPath, useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  session: Account;
};

export function useClusterCreateScreen(props: Props) {
  const router = useRouter();
  const directoryPath = useDirectoryPath();

  const initial: ClusterWithItems = {
    id: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    name: "",
    slug: "",
    description: "",
    owner: props.session,
    properties: {},
    items: [],
    clusters: [],
    assets: [],
  };

  async function handleCreate(cluster: ClusterInitialProps) {
    const created = await clusterCreate({
      name: cluster.name,
      slug: cluster.slug,
      url: cluster.url,
      description: cluster.description,
      content: cluster.content,
      asset_ids: cluster.asset_ids,
    });

    const newPath = joinDirectoryPath(directoryPath, created.slug);

    router.push(`/directory/${newPath}`);
  }

  return {
    initial,
    handlers: {
      handleCreate,
    },
  };
}
