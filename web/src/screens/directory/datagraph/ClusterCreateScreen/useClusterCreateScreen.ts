import { useRouter } from "next/navigation";

import { clusterCreate, clusterGet } from "src/api/openapi/clusters";
import { ClusterInitialProps, ClusterWithItems } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

import { joinDirectoryPath, useDirectoryPath } from "../useDirectoryPath";

export function useClusterCreateScreen() {
  const account = useSession();
  const router = useRouter();
  const directoryPath = useDirectoryPath();

  if (!account) {
    router.push("/login");
  }

  const initial: ClusterWithItems = {
    id: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    name: "",
    slug: "",
    description: "",
    owner: account!,
    properties: {},
    items: [],
    clusters: [],
    assets: [],
  };

  console.log({ directoryPath });

  async function handleCreate(cluster: ClusterInitialProps) {
    const created = await clusterCreate(cluster);

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
