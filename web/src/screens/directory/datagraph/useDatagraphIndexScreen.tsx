import { useClusterList } from "src/api/openapi/clusters";
import { useLinkList } from "src/api/openapi/links";
import {
  ClusterListOKResponse,
  LinkListOKResponse,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export type Props = {
  clusters: ClusterListOKResponse;
  links: LinkListOKResponse;
};

export function useDatagraphIndexScreen(props: Props) {
  const session = useSession();
  const {
    data: clusters,
    mutate: mutateClusters,
    error: errorClusters,
  } = useClusterList(
    {},
    {
      swr: {
        fallbackData: props.clusters,
      },
    },
  );

  const {
    data: links,
    mutate: mutateLinks,
    error: errorLinks,
  } = useLinkList(
    {},
    {
      swr: {
        fallbackData: props.links,
      },
    },
  );

  if (!clusters || !links) {
    return {
      ready: false as const,
      error: errorClusters || errorLinks,
    };
  }

  const empty = clusters.clusters.length === 0 && links.results === 0;

  return {
    ready: true as const,
    empty,
    data: {
      clusters: {
        data: clusters,
        mutate: mutateClusters,
      },
      links: {
        data: links,
        mutate: mutateLinks,
      },
    },
    mutate: {
      mutateClusters,
    },
    session,
  };
}
