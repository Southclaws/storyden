import { useClusterGet } from "src/api/openapi/clusters";
import { ClusterWithItems } from "src/api/openapi/schemas";

import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  slug: string;
  cluster: ClusterWithItems;
};

export function useClusterScreen(props: Props) {
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

  return {
    ready: true as const,
    data,
    directoryPath,
    mutate,
  };
}
