import { useSWRConfig } from "swr";

import { getThreadListKey } from "@/api/openapi-client/threads";

export function useFeedMutations() {
  const { mutate } = useSWRConfig();

  const threadQueryMutationKey = getThreadListKey()[0];

  return async () => {
    await mutate(
      (key) => Array.isArray(key) && key[0].startsWith(threadQueryMutationKey),
    );
  };
}
