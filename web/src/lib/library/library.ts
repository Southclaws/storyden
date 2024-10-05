import { MutatorCallback, useSWRConfig } from "swr";

import { getNodeListKey } from "@/api/openapi-client/nodes";
import { NodeListOKResponse, NodeListParams } from "@/api/openapi-schema";

export function useLibraryMutation(params?: NodeListParams) {
  const { mutate } = useSWRConfig();

  const nodeListKey = getNodeListKey(params);

  const revalidate = async (data?: MutatorCallback<NodeListOKResponse>) => {
    await mutate(
      (key) => Array.isArray(key) && key[0].startsWith(nodeListKey),
      data,
    );
  };

  return {
    revalidate,
  };
}
