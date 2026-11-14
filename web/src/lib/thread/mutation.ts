import { uniqueId } from "lodash";
import { useEffect, useRef } from "react";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import {
  postDelete,
  postReactAdd,
  postReactRemove,
  postUpdate,
} from "@/api/openapi-client/posts";
import { replyCreate } from "@/api/openapi-client/replies";
import { getThreadGetKey, threadUpdate } from "@/api/openapi-client/threads";
import {
  Account,
  Identifier,
  PostMutableProps,
  React,
  Reply,
  ReplyInitialProps,
  Thread,
  ThreadGetParams,
  ThreadGetResponse,
  ThreadReference,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";

export function useThreadMutations(
  thread: ThreadReference | ThreadGetResponse,
  currentPage?: number,
  totalPages?: number,
  initialSession?: Account,
) {
  // SWR renders fallbackData without placing it in the shared cache.
  // Reactions can start during that window, or when the client refresh fails.
  let initialThread: ThreadGetResponse | undefined;

  if ("replies" in thread) {
    initialThread = thread;
  }

  const sessionInitial = useSession(initialSession);
  const sessionRef = useRef(sessionInitial);
  useEffect(() => {
    sessionRef.current = sessionInitial;
  }, [sessionInitial]);

  const { mutate } = useSWRConfig();

  const threadGetKey = getThreadGetKey(thread.slug, {
    page: (currentPage ?? 1).toString(),
  });
  const key = (key: Arguments) => {
    if (!Array.isArray(key)) {
      return false;
    }

    const path = key[0];
    const params = key.length > 1 ? (key[1] as ThreadGetParams) : undefined;

    const pathMatch = path === threadGetKey[0];

    if (!pathMatch) {
      return false;
    }

    // Default to page 1 for comparison.
    const paramsMatch =
      (params?.page ?? "1") === (threadGetKey[1]?.page ?? "1");

    const match = pathMatch && paramsMatch;

    return match;
  };

  const createReply = async (reply: ReplyInitialProps) => {
    // Only apply optimistic update if we're on the last page
    // where the new reply will actually appear
    const isLastPage =
      !currentPage || !totalPages || currentPage === totalPages;

    if (isLastPage) {
      const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
        const session = sessionRef.current;

        if (!data || !session) {
          return;
        }

        const newReply = {
          id: uniqueId("optimistic_reply_"),
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          title: thread.title,
          author: session,
          assets: [],
          collections: { has_collected: false, in_collections: 0 },
          likes: { likes: 0, liked: false },
          reacts: [],
          body_links: [],
          slug: thread.slug,
          root_id: thread.id,
          root_slug: thread.slug,
          visibility: "published",
          ...reply,
          reply_to: undefined,
        } satisfies Reply;

        const newData: Thread = {
          ...data,
          replies: {
            ...data.replies,
            replies: [...data.replies.replies, newReply],
          },
        };

        return newData;
      };

      await mutate(key, mutator, {
        revalidate: false,
      });
    }

    return await replyCreate(thread.slug, reply);
  };

  const updateReply = async (id: Identifier, updated: PostMutableProps) => {
    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      if (!data) {
        return;
      }

      const newData = {
        ...data,
        replies: {
          ...data.replies,
          replies: data.replies.replies.map((reply) =>
            reply.id === id ? { ...reply, ...updated } : reply,
          ),
        },
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await postUpdate(id, updated);
  };

  const deleteReply = async (id: Identifier) => {
    const mutator = (data?: ThreadGetResponse) => {
      if (!data) {
        return;
      }

      const newData: ThreadGetResponse = {
        ...data,
        replies: {
          ...data.replies,
          replies: data.replies.replies.filter((reply) => reply.id !== id),
        },
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await postDelete(id);
  };

  const reactionAdd = async (postID: Identifier, emoji: string) => {
    const session = sessionRef.current;

    if (!session) {
      return;
    }

    const optimisticReaction = {
      id: uniqueId("optimistic_reply_update_"),
      emoji,
      author: session,
    } satisfies React;

    await mutate(
      key,
      (data: ThreadGetResponse | undefined = initialThread) =>
        updatePostReactions(data, postID, (reacts) => [
          ...reacts,
          optimisticReaction,
        ]),
      { revalidate: false },
    );

    try {
      const reaction = await postReactAdd(postID, { emoji });

      await mutate(
        key,
        (data: ThreadGetResponse | undefined = initialThread) =>
          updatePostReactions(data, postID, (reacts) =>
            reacts.map((item) => {
              if (item.id === optimisticReaction.id) {
                return reaction;
              }

              return item;
            }),
          ),
        { revalidate: false },
      );

      return reaction;
    } catch (error) {
      await mutate(
        key,
        (data: ThreadGetResponse | undefined = initialThread) =>
          updatePostReactions(data, postID, (reacts) =>
            reacts.filter((item) => item.id !== optimisticReaction.id),
          ),
        { revalidate: false },
      );

      throw error;
    }
  };

  const reactionRemove = async (postID: Identifier, reactID: Identifier) => {
    if (reactID.startsWith("optimistic")) {
      return;
    }

    const session = sessionRef.current;

    if (!session) {
      return;
    }

    let removedReaction: React | undefined;
    let removedIndex = -1;

    await mutate(
      key,
      (data: ThreadGetResponse | undefined = initialThread) =>
        updatePostReactions(data, postID, (reacts) => {
          const index = reacts.findIndex((reaction) => reaction.id === reactID);

          if (index === -1) {
            return reacts;
          }

          removedReaction = reacts[index];
          removedIndex = index;

          return reacts.filter((reaction) => reaction.id !== reactID);
        }),
      { revalidate: false },
    );

    try {
      await postReactRemove(postID, reactID);
    } catch (error) {
      const reaction = removedReaction;

      if (reaction) {
        await mutate(
          key,
          (data: ThreadGetResponse | undefined = initialThread) =>
            updatePostReactions(data, postID, (reacts) => {
              if (reacts.some((item) => item.id === reaction.id)) {
                return reacts;
              }

              const restored = [...reacts];
              const index = Math.max(
                0,
                Math.min(removedIndex, restored.length),
              );
              restored.splice(index, 0, reaction);

              return restored;
            }),
          { revalidate: false },
        );
      }

      throw error;
    }
  };

  const updateCategory = async (categoryID: Identifier) => {
    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      if (!data) {
        return;
      }

      const newData = {
        ...data,
        category_id: categoryID,
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await threadUpdate(thread.slug, { category: categoryID });
  };

  const revalidate = async (after?: number) => {
    setTimeout(() => {
      mutate(key);
    }, after);
  };

  return {
    createReply,
    updateReply,
    deleteReply,
    reactionAdd,
    reactionRemove,
    updateCategory,
    revalidate,
  };
}

function updatePostReactions(
  data: ThreadGetResponse | undefined,
  postID: Identifier,
  update: (reacts: React[]) => React[],
): ThreadGetResponse | undefined {
  if (!data) {
    return;
  }

  if (data.id === postID) {
    return { ...data, reacts: update(data.reacts) };
  }

  return {
    ...data,
    replies: {
      ...data.replies,
      replies: data.replies.replies.map((reply) => {
        if (reply.id !== postID) {
          return reply;
        }

        return { ...reply, reacts: update(reply.reacts) };
      }),
    },
  };
}
