import { useNodeList } from "src/api/openapi-client/nodes";
import {
  Account,
  NodeListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi-schema";

export type Props = {
  session: Account;
  threads: ThreadListOKResponse;
  nodes: NodeListOKResponse;
};

export function useDraftListScreen({ session, threads, nodes }: Props) {
  const { data: nodesData, error: errorNodes } = useNodeList(
    {
      author: session.handle,
      visibility: [Visibility.draft],
    },
    { swr: { fallbackData: nodes } },
  );

  if (!nodesData) {
    return {
      ready: false as const,
      error: errorNodes,
    };
  }

  const empty = nodes.nodes.length === 0;

  return {
    ready: true as const,
    empty,
    data: {
      nodes: nodesData.nodes,
      threads: threads.threads,
    },

    session,
  };
}
