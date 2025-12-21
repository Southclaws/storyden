import { groupBy, toPairs } from "lodash";
import { useEffect, useRef } from "react";

import { handle } from "@/api/client";
import { Account, React, ReactList, Reply, Thread } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useThreadMutations } from "@/lib/thread/mutation";

export const REACTION_THROTTLE = 180;

export type Props = {
  initialSession?: Account;
  thread: Thread;
  reply: Reply;
  currentPage?: number;
};

export type ReactCount = {
  emoji: string;
  count: number;
  hasReacted: boolean;
  reactions: React[];
};

function groupReactions(
  session: Account | undefined,
  reacts: ReactList,
): ReactCount[] {
  const grouped = groupBy<React>(reacts, "emoji");
  const pairs = toPairs<React[]>(grouped);

  return pairs.map(
    ([key, value]) =>
      ({
        emoji: key,
        count: value.length,
        hasReacted: Boolean(
          value.find((react) => react.author?.id === session?.id),
        ),
        reactions: value,
      }) satisfies ReactCount,
  );
}

export function useReactionList({ initialSession, thread, reply, currentPage }: Props) {
  const session = useSession(initialSession);
  const { reactionAdd, reactionRemove, revalidate } = useThreadMutations(
    thread,
    currentPage,
  );

  const isLoggedIn = Boolean(session);

  const postReactions = useRef(reply.reacts);
  useEffect(() => {
    postReactions.current = reply.reacts;
  }, [reply.reacts]);

  const reacts = groupReactions(session, reply.reacts);

  const handleAdd = async (emoji: string) => {
    await handle(async () => {
      await reactionAdd(reply.id, emoji);
    });
  };

  const handleRemove = async (id: string) => {
    await handle(async () => {
      await reactionRemove(reply.id, id);
    });
  };

  const handleReactExisting = (emoji: string, retry?: boolean) => {
    const currentReactions = postReactions.current;
    const grouped = groupBy<React>(currentReactions, "emoji");
    const reactions = grouped[emoji];

    const existing = reactions?.find(
      (r) => r.author?.id === session?.id && r.emoji === emoji,
    );

    if (existing) {
      if (existing.id.startsWith("optimistic") && !retry) {
        // If the selected reaction is not yet hydrated from the server, set a
        // timeout to re-try the deletion, ensuring that it's post-revalidation.
        setTimeout(
          () => handleReactExisting(emoji, true),
          REACTION_THROTTLE * 2,
        );
      } else {
        handleRemove(existing.id);
      }
    } else {
      handleAdd(emoji);
    }
  };

  // Only difference for the picker: if the user has already reacted with the
  // same emoji, don't do anything instead of removing it. Same as Discord.
  const handleReactPicker = (emoji: string) => {
    const currentReactions = postReactions.current;
    const grouped = groupBy<React>(currentReactions, "emoji");
    const reactions = grouped[emoji];

    const existing = reactions?.find(
      (r) => r.author?.id === session?.id && r.emoji === emoji,
    );

    if (existing) {
      return;
    }

    handleAdd(emoji);
  };

  return {
    data: {
      isLoggedIn,
      reacts,
    },
    handlers: {
      handleReactExisting,
      handleReactPicker,
    },
  };
}
