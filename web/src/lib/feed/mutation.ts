import { uniqueId } from "lodash/fp";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";
import { SWRInfiniteKeyLoader, unstable_serialize } from "swr/infinite";

import { cleanQuery } from "@/api/common";
import {
  getThreadListKey,
  threadCreate,
  threadDelete,
} from "@/api/openapi-client/threads";
import {
  Account,
  Identifier,
  LinkReference,
  ThreadInitialProps,
  ThreadListOKResponse,
  ThreadListParams,
  ThreadReference,
} from "@/api/openapi-schema";

type QueryParams = Record<string, string | string[] | undefined>;

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

export function useFeedMutations(session?: Account, params?: ThreadListParams) {
  const { mutate } = useSWRConfig();

  const pageKeyParams = {
    ...(params ? { categories: params.categories } : {}),
    ...(params?.page ? { page: params.page } : {}),
  } as ThreadListParams;

  const threadQueryMutationKey = unstable_serialize(
    getThreadListPageKey(pageKeyParams),
  );

  async function revalidate() {
    await mutate(threadQueryMutationKey);
  }

  async function createThread(
    initial: ThreadInitialProps,
    preHydratedLink?: LinkReference,
  ) {
    const mutator: MutatorCallback<ThreadListOKResponse[]> = (data) => {
      if (!data) return;
      if (!session) return;

      const description =
        new DOMParser()
          .parseFromString(initial.body, "text/html")
          .querySelector("body")?.textContent ?? "";

      const newThread = {
        ...initial,
        category: {
          id: initial.category,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          slug: "",
          admin: false,
          colour: "colour",
          description: "",
          name: "name",
          sort: 0,
        },
        id: uniqueId("optimistic_thread_id_"),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        slug: uniqueId("optimistic_thread_slug_"),
        author: session,
        description,
        body: initial.body,
        body_links: [],
        assets: [],
        collections: {
          has_collected: false,
          in_collections: 0,
        },
        likes: { likes: 0, liked: false },
        reacts: [],
        pinned: false,
        reply_status: { replies: 0, replied: 0 },
        tags: [],
        link: preHydratedLink,
      } satisfies ThreadReference;

      const newData = data.reduce((acc, page) => {
        // Append the new thread to the first page
        // NOTE: This assumes ordering is most recent first.
        if (page.current_page === 1) {
          acc.push({
            ...page,
            threads: [newThread, ...page.threads],
          });
        } else {
          acc.push(page);
        }

        return acc;
      }, [] as ThreadListOKResponse[]);

      return newData;
    };

    await mutate(threadQueryMutationKey, mutator, {
      revalidate: false,
    });

    return await threadCreate(initial);
  }

  async function deleteThread(id: Identifier) {
    const mutator: MutatorCallback<ThreadListOKResponse[]> = (data) => {
      if (!data) return;

      const newData = data.reduce((acc, page) => {
        // Scan every page, remove the deleted post.
        acc.push({
          ...page,
          threads: page.threads.filter((t) => t.id !== id),
        });

        return acc;
      }, [] as ThreadListOKResponse[]);

      return newData;
    };

    await mutate(threadQueryMutationKey, mutator, {
      revalidate: false,
    });

    await threadDelete(id);
  }

  return {
    createThread,
    deleteThread,
    revalidate,
  };
}
