import { useClusterList } from "src/api/openapi/clusters";
import {
  ClusterListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export type Props = {
  threads: ThreadListOKResponse;
  clusters: ClusterListOKResponse;
};

export function useDraftListScreen(props: Props) {
  const session = useSession();
  const {
    data: clusters,
    mutate: mutateClusters,
    error: errorClusters,
  } = useClusterList(
    {
      visibility: [Visibility.draft],
    },
    {
      swr: {
        fallbackData: props.clusters,
      },
    },
  );

  if (!clusters) {
    return {
      ready: false as const,
      error: errorClusters,
    };
  }

  const empty = clusters.clusters.length === 0;

  return {
    ready: true as const,
    empty,
    data: {
      clusters: {
        data: clusters,
        mutate: mutateClusters,
      },
    },
    mutate: {
      mutateClusters,
    },
    session,
  };
}
