import { uniqueId } from "lodash";
import { useEffect, useRef } from "react";
import { MutatorCallback, useSWRConfig } from "swr";

import {
  postDelete,
  postReactAdd,
  postReactRemove,
  postUpdate,
} from "@/api/openapi-client/posts";
import { replyCreate } from "@/api/openapi-client/replies";
import { getThreadGetKey } from "@/api/openapi-client/threads";
import {
  Identifier,
  PostMutableProps,
  React,
  Reply,
  ReplyInitialProps,
  Thread,
  ThreadGetResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";

export function useThreadMutations(thread: Thread) {
  const sessionInitial = useSession();
  const sessionRef = useRef(sessionInitial);
  useEffect(() => {
    sessionRef.current = sessionInitial;
  }, [sessionInitial]);

  const { mutate } = useSWRConfig();

  const key = getThreadGetKey(thread.slug);

  const createReply = async (reply: ReplyInitialProps) => {
    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      const session = sessionRef.current;

      if (!data || !session) return;

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
        ...reply,
      } satisfies Reply;

      const newData: Thread = {
        ...data,
        replies: [...data.replies, newReply],
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await replyCreate(thread.slug, reply);
  };

  const updateReply = async (id: Identifier, updated: PostMutableProps) => {
    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      if (!data) return;

      const newData = {
        ...data,
        replies: data.replies.map((reply) =>
          reply.id === id ? { ...reply, ...updated } : reply,
        ),
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
      if (!data) return;

      const newData: ThreadGetResponse = {
        ...data,
        replies: data.replies.filter((reply) => reply.id !== id),
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await postDelete(id);
  };

  const reactionAdd = async (replyID: Identifier, emoji: string) => {
    const session = sessionRef.current;
    if (!session) return;

    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      if (!data) return;

      const newData = {
        ...data,
        replies: data.replies.map((reply) => {
          if (reply.id !== replyID) {
            return reply;
          }

          const newReact = {
            id: uniqueId("optimistic_reply_update_"),
            emoji,
            author: session,
          } satisfies React;

          const reacts = [...reply.reacts, newReact];

          return {
            ...reply,
            reacts,
          };
        }),
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await postReactAdd(replyID, { emoji });
  };

  const reactionRemove = async (replyID: Identifier, reactID: Identifier) => {
    const session = sessionRef.current;
    if (!session) return;

    const mutator: MutatorCallback<ThreadGetResponse> = (data) => {
      if (!data) return;

      const newData = {
        ...data,
        replies: data.replies.map((reply) => {
          if (reply.id !== replyID) {
            return reply;
          }

          const reacts = reply.reacts.filter((react) => react.id !== reactID);

          return {
            ...reply,
            reacts,
          };
        }),
      };

      return newData;
    };

    await mutate(key, mutator, {
      revalidate: false,
    });

    await postReactRemove(replyID, reactID);
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
    revalidate,
  };
}
