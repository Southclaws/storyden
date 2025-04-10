import { PostReference } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useFeedMutations } from "@/lib/feed/mutation";

export type Props = {
  thread: PostReference;
};

export function useLikeButton({ thread }: Props) {
  const { likePost, unlikePost, revalidate } = useFeedMutations();

  const handleClick = async () => {
    // This is a toggle button, so we'll either like or unlike depending on what they have now.
    if (thread.likes.liked) {
        await handle(
          async () => {
            await unlikePost(thread.id);
          },
          {
            cleanup: async () => await revalidate(),
          },
        );
    } else {
      await handle(
        async () => {
          await likePost(thread.id);
        },
        {
          cleanup: async () => await revalidate(),
        },
      );
    }
  }

  return {
    ready: true as const,
    handleClick
  };
}
