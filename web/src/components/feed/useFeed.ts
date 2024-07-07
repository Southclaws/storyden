import { filter } from "lodash/fp";
import { useSWRConfig } from "swr";

import { useNodeList } from "src/api/openapi/nodes";
import {
  NodeList,
  NodeListParams,
  NodeListResult,
  ThreadList,
  ThreadListParams,
  ThreadListResult,
  ThreadMark,
  ThreadReference,
} from "src/api/openapi/schemas";
import {
  getThreadListKey,
  threadDelete,
  useThreadList,
} from "src/api/openapi/threads";

export type MixedContent = {
  threads: ThreadListResult;
  nodes?: NodeListResult;
};

export type MixedContentLists = {
  threads: ThreadList;
  nodes: NodeList;
};

export type MixedContentHandlers = {
  handleDeleteThread: (id: string) => void;
};

export type Props = {
  params?: {
    threads?: ThreadListParams;
    nodes?: NodeListParams;
  };
  initialData?: {
    threads: ThreadListResult;
    nodes?: NodeListResult;
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
    data: nodes,
    mutate: mutateNodes,
    error: errorNodes,
  } = useNodeList(params?.nodes, {
    swr: { fallbackData: initialData?.nodes },
  });

  const isReady = threads && nodes;
  const allErrors = errorThreads || errorNodes;

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
      nodes: initialData?.nodes ? nodes : undefined,
    } satisfies MixedContent,
    mutate: {
      mutateThreads,
      mutateNodes,
    },
    handlers: {
      handleDeleteThread,
    } satisfies MixedContentHandlers,
  };
}

export function useFeedMutation() {
  const { mutate } = useSWRConfig();

  const threadQueryMutationKey = getThreadListKey()[0];

  return async () => {
    await mutate(
      (key) => Array.isArray(key) && key[0].startsWith(threadQueryMutationKey),
    );
  };
}
