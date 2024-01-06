import { useClusterList } from "src/api/openapi/clusters";
import { useItemList } from "src/api/openapi/items";
import {
  ClusterListOKResponse,
  ItemListOKResponse,
} from "src/api/openapi/schemas";

export type Props = {
  clusters: ClusterListOKResponse;
  items: ItemListOKResponse;
};

export function useDatagraphIndexScreen(props: Props) {
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
    data: items,
    mutate: mutateItems,
    error: errorItems,
  } = useItemList(
    {},
    {
      swr: {
        fallbackData: props.items,
      },
    },
  );

  if (!clusters || !items) {
    return {
      ready: false as const,
      error: errorClusters || errorItems,
    };
  }

  return {
    ready: true as const,
    data: {
      clusters: {
        data: clusters,
        mutate: mutateClusters,
      },
      items: {
        data: items,
        mutate: mutateItems,
      },
    },
    mutate: {
      mutateClusters,
      mutateItems,
    },
  };
}
