"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { handle } from "@/api/client";
import { PostReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useConfirmation } from "@/components/site/useConfirmation";
import { WEB_ADDRESS } from "@/config";
import { useFeedMutations } from "@/lib/feed/mutation";
import { canDeletePost } from "@/lib/thread/permissions";
import { useShare } from "@/utils/client";

export type Props = {
  thread: PostReference;
};

export function useFeedItemMenu({ thread }: Props) {
  const account = useSession();
  const permalink = getPermalinkForThread(thread.slug);
  const [, copyToClipboard] = useCopyToClipboard();

  const { deleteThread, revalidate } = useFeedMutations();

  const {
    isConfirming: isConfirmingDelete,
    handleConfirmAction: handleConfirmDelete,
    handleCancelAction: handleCancelDelete,
  } = useConfirmation(handleDelete);

  const isSharingEnabled = useShare();
  const isDeletingEnabled = canDeletePost(thread, account);

  async function share() {
    await navigator.share({
      title: `A post by ${thread.author.name}`,
      url: permalink,
      text: thread.description,
    });
  }

  async function handleDelete() {
    handle(async () => await deleteThread(thread.id), {
      cleanup: async () => await revalidate(),
    });
  }

  function handleSelect({ value }: { value: string }) {
    switch (value) {
      case "copy-link":
        copyToClipboard(permalink);
        return;

      case "share":
        share();
        return;

      case "delete":
        handleConfirmDelete();
        return;

      default:
        throw new Error("unknown handler");
    }
  }

  return {
    isSharingEnabled,
    isDeletingEnabled,
    isConfirmingDelete,
    handlers: {
      handleSelect,
      handleCancelDelete,
    },
  };
}

function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}
