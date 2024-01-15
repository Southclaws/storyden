import { useRouter } from "next/navigation";

import { clusterUpdate, useClusterGet } from "src/api/openapi/clusters";
import { ClusterInitialProps, ClusterWithItems } from "src/api/openapi/schemas";

import { replaceDirectoryPath, useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  cluster: ClusterWithItems;
};

export function useClusterViewerScreen(props: Props) {
  const router = useRouter();
  const { data, mutate, error } = useClusterGet(props.slug, {
    swr: {
      fallbackData: props.cluster,
    },
  });

  const directoryPath = useDirectoryPath();

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const { slug } = data;

  async function handleSave(cluster: ClusterInitialProps) {
    await clusterUpdate(slug, {
      name: cluster.name,
      slug: cluster.slug,
      asset_ids: cluster.asset_ids,
      url: cluster.url,
      description: cluster.description,
      content: cluster.content,
      properties: cluster.properties,
    });
    await mutate();

    // Handle slug changes properly by redirecting to the new path.
    if (cluster.slug !== slug) {
      const newPath = replaceDirectoryPath(directoryPath, slug, cluster.slug);
      router.push(newPath);
    }
  }

  return {
    ready: true as const,
    data,
    handlers: { handleSave },
    directoryPath,
    mutate,
  };
}
