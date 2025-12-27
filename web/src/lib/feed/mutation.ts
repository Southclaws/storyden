import { uniqueId } from "lodash/fp";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import { likePostAdd, likePostRemove } from "@/api/openapi-client/likes";
import {
  getThreadListKey,
  threadCreate,
  threadDelete,
  threadUpdate,
} from "@/api/openapi-client/threads";
import {
  Account,
  Identifier,
  LinkReference,
  ThreadInitialProps,
  ThreadListOKResponse,
  ThreadListParams,
  ThreadMutableProps,
  ThreadReference,
  Visibility,
} from "@/api/openapi-schema";
import { deepEqual } from "@/utils/equality";

export function useFeedMutations(session?: Account, params?: ThreadListParams) {
  const { mutate } = useSWRConfig();

  const pageKeyParams = {
    ...(params ? { categories: params.categories } : {}),
    ...(params?.page ? { page: params.page } : {}),
  } as ThreadListParams;

  const threadQueryMutationKey = getThreadListKey(pageKeyParams);
  function threadListKeyFilterFn(key: Arguments) {
    if (!Array.isArray(key)) return false;

    const path = key[0];
    const params = key[1] as ThreadListParams;

    const pathMatch = path === threadQueryMutationKey[0];
    if (!pathMatch) return false;

    const pageMatch = (params.page ?? "1") === (pageKeyParams.page ?? "1");

    const categoryMatch = pageKeyParams.categories
      ? deepEqual(params.categories, pageKeyParams.categories)
      : true;

    const paramsMatch = pageMatch && categoryMatch;

    return pathMatch && paramsMatch;
  }

  async function revalidate() {
    await mutate(threadListKeyFilterFn);
  }

  async function createThread(
    initial: ThreadInitialProps,
    preHydratedLink?: LinkReference,
  ) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;
      if (!session) return;

      // Don't mutate if not on first page.
      if (data.current_page === 1) {
        return;
      }

      const description =
        new DOMParser()
          .parseFromString(initial.body ?? "<body></body>", "text/html")
          .querySelector("body")?.textContent ?? "";

      const newThread = {
        ...initial,
        category: initial.category
          ? {
              id: initial.category,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
              slug: "",
              colour: "colour",
              description: "",
              name: "name",
              sort: 0,
              children: [],
            }
          : undefined,
        id: uniqueId("optimistic_thread_id_"),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        slug: uniqueId("optimistic_thread_slug_"),
        author: session,
        description,
        body: initial.body ?? "",
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
        visibility: Visibility.draft,
        link: preHydratedLink,
      } satisfies ThreadReference;

      const newData = {
        ...data,
        threads: [newThread, ...data.threads],
      };

      return newData;
    };

    await mutate(threadListKeyFilterFn, mutator, {
      revalidate: false,
    });

    return await threadCreate(initial);
  }

  async function updateThread(id: Identifier, updated: ThreadMutableProps) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      return {
        ...data,
        threads: data.threads.map((t) => {
          if (t.id !== id) {
            return t;
          }

          const newThread = {
            ...t,
            ...updated,
            category: t.category
              ? {
                  ...t.category,
                  id: updated.category ?? t.category.id,
                }
              : updated.category
                ? {
                    id: updated.category,
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString(),
                    slug: "",
                    colour: "",
                    description: "",
                    name: "",
                    sort: 0,
                    children: [],
                  }
                : undefined,
            tags: t.tags,
          } satisfies ThreadReference;

          return newThread;
        }),
      };
    };

    await mutate(threadListKeyFilterFn, mutator, {
      revalidate: false,
    });

    await threadUpdate(id, updated);
  }

  async function deleteThread(id: Identifier) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      return {
        ...data,
        threads: data.threads.filter((t) => t.id !== id),
      };
    };

    await mutate(threadListKeyFilterFn, mutator, {
      revalidate: false,
    });

    await threadDelete(id);
  }

  async function likePost(id: Identifier) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      const newThreads = data.threads.map((thread) => {
        if (thread.id === id) {
          return {
            ...thread,
            likes: {
              likes: thread.likes.likes + 1,
              liked: true,
            },
          };
        }
        return thread;
      });

      return {
        ...data,
        threads: newThreads,
      };
    };

    await mutate(threadListKeyFilterFn, mutator, {
      revalidate: false,
    });

    await likePostAdd(id);
  }

  async function unlikePost(id: Identifier) {
    const mutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      const newThreads = data.threads.map((thread) => {
        if (thread.id === id) {
          return {
            ...thread,
            likes: {
              likes: thread.likes.likes - 1,
              liked: false,
            },
          };
        }
        return thread;
      });

      return {
        ...data,
        threads: newThreads,
      };
    };

    await mutate(threadListKeyFilterFn, mutator, {
      revalidate: false,
    });

    await likePostRemove(id);
  }

  return {
    createThread,
    updateThread,
    deleteThread,
    likePost,
    unlikePost,
    revalidate,
  };
}
