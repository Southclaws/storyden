import { useClusterList } from "src/api/openapi/clusters";
import { useItemList } from "src/api/openapi/items";
import { useLinkList } from "src/api/openapi/links";
import {
  ClusterListOKResponse,
  ItemListOKResponse,
  LinkListOKResponse,
} from "src/api/openapi/schemas";

export type Props = {
  clusters: ClusterListOKResponse;
  items: ItemListOKResponse;
  links: LinkListOKResponse;
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

  if (!clusters || !items || !links) {
    return {
      ready: false as const,
      error: errorClusters || errorItems || errorLinks,
    };
  }

  const empty =
    clusters.clusters.length === 0 &&
    items.items.length === 0 &&
    links.results === 0;

  return {
    ready: true as const,
    empty,
    data: {
      clusters: {
        data: clusters,
        mutate: mutateClusters,
      },
      items: {
        data: items,
        mutate: mutateItems,
      },
      links: {
        data: links,
        mutate: mutateLinks,
      },
    },
    mutate: {
      mutateClusters,
      mutateItems,
    },
  };
}
