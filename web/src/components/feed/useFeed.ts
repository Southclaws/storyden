import { filter } from "lodash/fp";

import { useLinkList } from "src/api/openapi/links";
import { useNodeList } from "src/api/openapi/nodes";
import {
  LinkList,
  LinkListParams,
  LinkListResult,
  NodeList,
  NodeListParams,
  NodeListResult,
  ThreadList,
  ThreadListParams,
  ThreadListResult,
  ThreadMark,
  ThreadReference,
} from "src/api/openapi/schemas";
import { threadDelete, useThreadList } from "src/api/openapi/threads";

export type MixedContent = {
  threads: ThreadListResult;
  nodes: NodeListResult;
  links: LinkListResult;
};

export type MixedContentLists = {
  threads: ThreadList;
  nodes: NodeList;
  links: LinkList;
};

export type MixedContentHandlers = {
  handleDeleteThread: (id: string) => void;
};

export type Props = {
  params?: {
    threads?: ThreadListParams;
    nodes?: NodeListParams;
    links?: LinkListParams;
  };
  initialData?: {
    threads: ThreadListResult;
    nodes: NodeListResult;
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
    data: nodes,
    mutate: mutateNodes,
    error: errorNodes,
  } = useNodeList(params?.nodes, {
    swr: { fallbackData: initialData?.nodes },
  });

  const {
    data: links,
    mutate: mutateLinks,
    error: errorLinks,
  } = useLinkList(params?.links, {
    swr: { fallbackData: initialData?.links },
  });

  const isReady = threads && nodes && links;
  const allErrors = errorThreads || errorNodes || errorLinks;

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
      nodes,
      links,
    } satisfies MixedContent,
    mutate: {
      mutateThreads,
      mutateNodes,
      mutateLinks,
    },
    handlers: {
      handleDeleteThread,
    } satisfies MixedContentHandlers,
  };
}
