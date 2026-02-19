import { PostReference } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useFeedMutations } from "@/lib/feed/mutation";

export type Props = {
  thread: PostReference;
};

export function useLikeButton({ thread }: Props) {
  const session = useSession();
  const { likePost, unlikePost, revalidate } = useFeedMutations();

  const enabled = Boolean(session);

  const handleClick = async () => {
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
          await revalidate();
        },
      },
    );
  };

  return {
    handleClick,
    enabled,
  };
}
