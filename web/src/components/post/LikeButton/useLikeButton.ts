import { useState } from "react";

import { handle } from "@/api/client";
import { Account, PostReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useFeedMutations } from "@/lib/feed/mutation";

export type Props = {
  initialSession?: Account;
  thread: PostReference;
};

export function useLikeButton({ initialSession, thread }: Props) {
  const session = useSession(initialSession);
  const { likePost, unlikePost, revalidate } = useFeedMutations();

  const enabled = Boolean(session);
  const [pending, setPending] = useState(false);

  const handleClick = async () => {
    if (!enabled || pending) return;
    setPending(true);
    await handle(
      async () => {
        // This is a toggle button, so we'll either like or unlike depending on what they have now.
        if (thread.likes.liked) {
          await unlikePost(thread.id);
        } else {
          await likePost(thread.id);
        }
      },
      {
        async cleanup() {
          setPending(false);
          await revalidate();
        },
      },
    );
  };

  return {
    handleClick,
    enabled,
    pending,
  };
}
