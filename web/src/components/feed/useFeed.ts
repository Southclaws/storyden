import { filter } from "lodash/fp";

import { useClusterList } from "src/api/openapi/clusters";
import { useItemList } from "src/api/openapi/items";
import { useLinkList } from "src/api/openapi/links";
import {
  ClusterList,
  ClusterListParams,
  ClusterListResult,
  ItemList,
  ItemListParams,
  ItemListResult,
  LinkList,
  LinkListParams,
  LinkListResult,
  ThreadList,
  ThreadListParams,
  ThreadListResult,
  ThreadMark,
  ThreadReference,
} from "src/api/openapi/schemas";
import { threadDelete, useThreadList } from "src/api/openapi/threads";

export type MixedContent = {
  threads: ThreadListResult;
  clusters: ClusterListResult;
  items: ItemListResult;
  links: LinkListResult;
};

export type MixedContentLists = {
  threads: ThreadList;
  clusters: ClusterList;
  items: ItemList;
  links: LinkList;
};

export type MixedContentHandlers = {
  handleDeleteThread: (id: string) => void;
};

export type Props = {
  params?: {
    threads?: ThreadListParams;
    clusters?: ClusterListParams;
    items?: ItemListParams;
    links?: LinkListParams;
  };
  initialData?: {
    threads: ThreadListResult;
    clusters: ClusterListResult;
    items: ItemListResult;
    links: LinkListResult;
  };
};

const removeThread = (id: string) =>
  filter((v: ThreadReference) => v.id !== id);

export function useFeed({ params, initialData }: Props) {
  const {
    data: threads,
    mutate: mutateThreads,
    error: errorThreads,
  } = useThreadList(params?.threads, {
    swr: { fallbackData: initialData?.threads },
  });

  const {
    data: clusters,
    mutate: mutateClusters,
    error: errorClusters,
  } = useClusterList(params?.clusters, {
    swr: { fallbackData: initialData?.clusters },
  });

  const {
    data: items,
    mutate: mutateItems,
    error: errorItems,
  } = useItemList(params?.items, {
    swr: { fallbackData: initialData?.items },
  });

  const {
    data: links,
    mutate: mutateLinks,
    error: errorLinks,
  } = useLinkList(params?.links, {
    swr: { fallbackData: initialData?.links },
  });

  const isReady = threads && clusters && items && links;
  const allErrors = errorThreads || errorClusters || errorItems || errorLinks;

  if (!isReady) {
    return {
      ready: false as const,
      error: allErrors,
    };
  }

  async function handleDeleteThread(id: ThreadMark) {
    await threadDelete(id);

    const existingThreads =
      threads?.threads ?? initialData?.threads.threads ?? [];
    const newThreads = removeThread(id)(existingThreads);

    if (initialData) {
      mutateThreads({
        ...initialData?.threads,
        threads: newThreads,
      });
    } else {
      mutateThreads();
    }
  }

  return {
    ready: true as const,
    data: {
      threads,
      clusters,
      items,
      links,
    } satisfies MixedContent,
    mutate: {
      mutateThreads,
      mutateClusters,
      mutateItems,
      mutateLinks,
    },
    handlers: {
      handleDeleteThread,
    } satisfies MixedContentHandlers,
  };
}
