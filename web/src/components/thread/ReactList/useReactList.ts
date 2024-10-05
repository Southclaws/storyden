import { groupBy, toPairs } from "lodash";
import { useEffect, useRef } from "react";

import { handle } from "@/api/client";
import { Account, React, ReactList, Reply, Thread } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useThreadMutations } from "@/lib/thread/mutation";

export const REACTION_THROTTLE = 2500;

export type Props = {
  thread: Thread;
  reply: Reply;
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

export function useReactionList({ thread, reply }: Props) {
  const session = useSession();
  const { reactionAdd, reactionRemove, revalidate } =
    useThreadMutations(thread);

  const postReactions = useRef(reply.reacts);
  useEffect(() => {
    postReactions.current = reply.reacts;
  }, [reply.reacts]);

  const reacts = groupReactions(session, reply.reacts);

  const hasAlreadyReacted = reply.reacts.find(
    (react) => react.author?.id === session?.id,
  );

  const handleAdd = async (emoji: string) => {
    await handle(
      async () => {
        await reactionAdd(reply.id, emoji);
      },
      { cleanup: () => revalidate(REACTION_THROTTLE) },
    );
  };

  const handleRemove = async (id: string) => {
    await handle(
      async () => {
        await reactionRemove(reply.id, id);
      },
      { cleanup: () => revalidate(REACTION_THROTTLE) },
    );
  };

  const handleReactExisting = (emoji: string) => {
    const currentReactions = postReactions.current;
    const grouped = groupBy<React>(currentReactions, "emoji");
    const reactions = grouped[emoji];

    const existing = reactions?.find(
      (r) => r.author?.id === session?.id && r.emoji === emoji,
    );

    if (existing) {
      handleRemove(existing.id);
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
      hasAlreadyReacted,
      reacts,
    },
    handlers: {
      handleReactExisting,
      handleReactPicker,
    },
  };
}
