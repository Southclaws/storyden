import { filter } from "lodash/fp";
import { useSWRConfig } from "swr";

import {
  getThreadListKey,
  useThreadList,
} from "src/api/openapi-client/threads";
import {
  ThreadListParams,
  ThreadListResult,
  ThreadReference,
} from "src/api/openapi-schema";

export type Props = {
  params?: ThreadListParams;
  initialData?: ThreadListResult;
};

const removeThread = (id: string) =>
  filter((v: ThreadReference) => v.id !== id);

export function useFeed({ params, initialData }: Props) {
  const { data, mutate, error } = useThreadList(params, {
    swr: { fallbackData: initialData },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data,
    mutate,
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
