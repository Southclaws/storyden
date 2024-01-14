import { useRouter } from "next/navigation";

import { clusterCreate } from "src/api/openapi/clusters";
import { ClusterInitialProps, ClusterWithItems } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export function useClusterCreateScreen() {
  const account = useSession();
  const router = useRouter();

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
  };

  async function handleCreate(cluster: ClusterInitialProps) {
    await clusterCreate(cluster);
  }

  return {
    initial,
    handlers: {
      handleCreate,
    },
  };
}
