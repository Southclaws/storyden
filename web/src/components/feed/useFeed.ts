import { filter } from "lodash/fp";

import {
  ThreadList,
  ThreadListParams,
  ThreadMark,
  ThreadReference,
} from "src/api/openapi/schemas";
import { threadDelete, useThreadList } from "src/api/openapi/threads";

const removeThread = (id: string) =>
  filter((v: ThreadReference) => v.id !== id);

export function useFeed(
  params?: ThreadListParams,
  initialThreads?: ThreadList,
) {
  const { data, error, mutate } = useThreadList(params, {
    swr: {
      fallbackData: initialThreads && { threads: initialThreads },
    },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  async function handleDelete(id: ThreadMark) {
    await threadDelete(id);

    const existingThreads = data?.threads ?? initialThreads ?? [];
    const newThreads = removeThread(id)(existingThreads);

    console.log({ existingThreads, newThreads });

    mutate({
      threads: newThreads,
    });
  }

  return {
    ready: true as const,
    data,
    handlers: {
      handleDelete,
    },
  };
}
