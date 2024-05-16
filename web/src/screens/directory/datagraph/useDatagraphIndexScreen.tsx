import { useLinkList } from "src/api/openapi/links";
import { useNodeList } from "src/api/openapi/nodes";
import {
  LinkListOKResponse,
  NodeListOKResponse,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export type Props = {
  nodes: NodeListOKResponse;
  links: LinkListOKResponse;
};

export function useDatagraphIndexScreen(props: Props) {
  const session = useSession();
  const {
    data: nodes,
    mutate: mutateNodes,
    error: errorNodes,
  } = useNodeList(
    {},
    {
      swr: {
        fallbackData: props.nodes,
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

  if (!nodes || !links) {
    return {
      ready: false as const,
      error: errorNodes || errorLinks,
    };
  }

  const empty = nodes.nodes.length === 0 && links.results === 0;

  return {
    ready: true as const,
    empty,
    data: {
      nodes: {
        data: nodes,
        mutate: mutateNodes,
      },
      links: {
        data: links,
        mutate: mutateLinks,
      },
    },
    mutate: {
      mutateNodes,
    },
    session,
  };
}
