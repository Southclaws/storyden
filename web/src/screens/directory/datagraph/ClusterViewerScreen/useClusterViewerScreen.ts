import { useRouter } from "next/navigation";

import {
  clusterDelete,
  clusterUpdate,
  clusterUpdateVisibility,
  useClusterGet,
} from "src/api/openapi/clusters";
import {
  Cluster,
  ClusterInitialProps,
  ClusterWithChildren,
  Visibility,
} from "src/api/openapi/schemas";

import { replaceDirectoryPath } from "../directory-path";
import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  cluster: ClusterWithChildren;
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

  async function handleVisibilityChange(visibility: Visibility) {
    await clusterUpdateVisibility(slug, { visibility });
    await mutate();
  }

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

  // TODO: Provide a way to set the new parent cluster for child clusters/items.
  async function handleDelete(cluster: Cluster) {
    const { destination } = await clusterDelete(cluster.slug);

    if (destination) {
      const newPath = replaceDirectoryPath(
        directoryPath,
        slug,
        destination.slug,
      );
      router.push(newPath);
    } else {
      router.push("/directory");
    }
  }

  return {
    ready: true as const,
    data,
    handlers: { handleSave, handleVisibilityChange, handleDelete },
    directoryPath,
    mutate,
  };
}
