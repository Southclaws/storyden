import { useNodeList } from "src/api/openapi-client/nodes";
import { useThreadList } from "src/api/openapi-client/threads";
import {
  Account,
  NodeListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi-schema";

export type Props = {
  session: Account;
  initialThreads: ThreadListOKResponse;
  initialNodes: NodeListOKResponse;
};

export function useDraftListScreen({
  session,
  initialThreads,
  initialNodes,
}: Props) {
  const { data: threadsData, error: errorThreads } = useThreadList(
    {
      author: session.handle,
      visibility: [Visibility.draft],
    },
    { swr: { fallbackData: initialThreads } },
  );

  const { data: nodesData, error: errorNodes } = useNodeList(
    {
      author: session.handle,
      visibility: [Visibility.draft],
    },
    { swr: { fallbackData: initialNodes } },
  );

  if (!threadsData || !nodesData) {
    return {
      ready: false as const,
      error: errorThreads || errorNodes,
    };
  }

  const empty =
    initialNodes.nodes.length === 0 && initialThreads.threads.length === 0;

  return {
    ready: true as const,
    empty,
    data: {
      nodes: nodesData.nodes,
      threads: threadsData.threads,
    },

    session,
  };
}
