import { useNodeList } from "src/api/openapi-client/nodes";
import {
  NodeListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi-schema";
import { useSession } from "src/auth";

export type Props = {
  threads: ThreadListOKResponse;
  nodes: NodeListOKResponse;
};

export function useDraftListScreen(props: Props) {
  const session = useSession();
  const {
    data: nodes,
    mutate: mutateNodes,
    error: errorNodes,
  } = useNodeList(
    {
      visibility: [Visibility.draft],
    },
    {
      swr: {
        fallbackData: props.nodes,
      },
    },
  );

  if (!nodes) {
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
      nodes: {
        data: nodes,
        mutate: mutateNodes,
      },
    },
    mutate: {
      mutateNodes,
    },
    session,
  };
}
