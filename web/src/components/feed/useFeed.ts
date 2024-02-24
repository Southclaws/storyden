import { filter } from "lodash/fp";

import {
  ClusterListParams,
  ClusterListResult,
  ItemListParams,
  ItemListResult,
  LinkListParams,
  LinkListResult,
  ThreadListParams,
  ThreadListResult,
  ThreadMark,
  ThreadReference,
} from "src/api/openapi/schemas";
import { threadDelete, useThreadList } from "src/api/openapi/threads";

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
  const { data, error, mutate } = useThreadList(params?.threads, {
    swr: { fallbackData: initialData?.threads },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  async function handleDeleteThread(id: ThreadMark) {
    await threadDelete(id);

    const existingThreads = data?.threads ?? initialData?.threads.threads ?? [];
    const newThreads = removeThread(id)(existingThreads);

    if (initialData) {
      mutate({
        ...initialData?.threads,
        threads: newThreads,
      });
    } else {
      mutate();
    }
  }

  return {
    ready: true as const,
    data,
    mutate,
    handlers: {
      handleDeleteThread,
    },
  };
}
