import { useEffect, useRef, useState } from "react";

import { handle } from "@/api/client";
import { Account, Reply, Thread } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useThreadMutations } from "@/lib/thread/mutation";

import { groupReactions } from "./ReactList.model";

export type Props = {
  initialSession?: Account;
  thread: Thread;
  reply: Reply;
  currentPage?: number;
};

export function useReactionList({
  initialSession,
  thread,
  reply,
  currentPage,
}: Props) {
  const session = useSession(initialSession);
  const { reactionAdd, reactionRemove } = useThreadMutations(
    thread,
    currentPage,
  );
  const pendingRef = useRef(new Set<string>());
  const [pendingEmojis, setPendingEmojis] = useState<ReadonlySet<string>>(
    () => new Set(),
  );
  const postReactions = useRef(reply.reacts);

  useEffect(() => {
    postReactions.current = reply.reacts;
  }, [reply.reacts]);

  async function runReaction(emoji: string, operation: () => Promise<unknown>) {
    if (pendingRef.current.has(emoji)) {
      return;
    }

    pendingRef.current.add(emoji);
    setPendingEmojis(new Set(pendingRef.current));

    return operation().finally(() => {
      pendingRef.current.delete(emoji);
      setPendingEmojis(new Set(pendingRef.current));
    });
  }

  async function handleAdd(emoji: string) {
    await handle(async () => {
      await reactionAdd(reply.id, emoji);
    });
  }

  async function handleRemove(id: string) {
    await handle(async () => {
      await reactionRemove(reply.id, id);
    });
  }

  async function handleReactExisting(emoji: string) {
    if (!session || pendingRef.current.has(emoji)) {
      return;
    }

    const existing = postReactions.current.find(
      (reaction) =>
        reaction.emoji === emoji && reaction.author?.id === session.id,
    );

    if (existing?.id.startsWith("optimistic")) {
      return;
    }

    await runReaction(emoji, () =>
      existing ? handleRemove(existing.id) : handleAdd(emoji),
    );
  }

  async function handleReactPicker(emoji: string) {
    if (!session || pendingRef.current.has(emoji)) {
      return;
    }

    const existing = postReactions.current.some(
      (reaction) =>
        reaction.emoji === emoji && reaction.author?.id === session.id,
    );
    if (existing) {
      return;
    }

    await runReaction(emoji, () => handleAdd(emoji));
  }

  return {
    data: {
      isLoggedIn: Boolean(session),
      pendingEmojis,
      reactions: groupReactions(session?.id, reply.reacts),
    },
    handlers: {
      handleReactExisting,
      handleReactPicker,
    },
  };
}
