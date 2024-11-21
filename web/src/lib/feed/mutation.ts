import { Arguments, MutatorCallback, useSWRConfig } from "swr";
import { SWRInfiniteKeyLoader } from "swr/infinite";

import { cleanQuery } from "@/api/common";
import { getThreadListKey, threadDelete } from "@/api/openapi-client/threads";
import { Identifier, ThreadListOKResponse } from "@/api/openapi-schema";

type QueryParams = Record<string, string | undefined>;

export const getThreadListPageKey =
  (parameters?: QueryParams): SWRInfiniteKeyLoader<ThreadListOKResponse> =>
  (pageIndex: number, previousPageData: ThreadListOKResponse | null) => {
    if (previousPageData && previousPageData.next_page === undefined) {
      return null;
    }

    const pageNumber = pageIndex + 1;

    const [path, params] = getThreadListKey({
      page: pageNumber.toString(),
      ...parameters,
    });

    const key = path + cleanQuery(params);

    return key;
  };

export function useFeedMutations() {
  const { mutate } = useSWRConfig();

  const threadQueryMutationKey = getThreadListKey()[0];

  function keyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(threadQueryMutationKey);
  }

  async function revalidate() {
    await mutate(keyFilterFn);
  }

  async function deleteThread(id: Identifier) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data || !data.threads) return;

      const newData = {
        ...data,
        threads: data.threads.filter((t) => t.id !== id),
      };

      return newData;
    };

    await mutate(keyFilterFn, mutator, {
      revalidate: false,
    });

    await threadDelete(id);
  }

  return {
    deleteThread,
    revalidate,
  };
}
