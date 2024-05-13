import { last } from "lodash";
import { useRouter } from "next/navigation";

import { clusterCreate } from "src/api/openapi/clusters";
import {
  Account,
  ClusterInitialProps,
  ClusterWithChildren,
} from "src/api/openapi/schemas";

import { joinDirectoryPath } from "../directory-path";
import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  session: Account;
};

export function useClusterCreateScreen(props: Props) {
  const router = useRouter();
  const directoryPath = useDirectoryPath();

  const initial: ClusterWithChildren = {
    id: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    name: "",
    slug: "",
    description: "",
    owner: props.session,
    properties: {},
    children: [],
    assets: [],
    visibility: "draft",
    recomentations: [],
  };

  async function handleCreate(cluster: ClusterInitialProps) {
    const parentSlug = last(directoryPath as string[]);
    const created = await clusterCreate({
      name: cluster.name,
      slug: cluster.slug,
      url: cluster.url,
      description: cluster.description,
      content: cluster.content,
      asset_ids: cluster.asset_ids,
      parent: parentSlug,
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
