"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";

import { ThreadReference } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useConfirmation } from "@/components/site/useConfirmation";
import { useFeedMutations } from "@/lib/feed/mutation";
import { canDeletePost, canEditPost } from "@/lib/thread/permissions";
import { withUndo } from "@/lib/thread/undo";
import { useShare } from "@/utils/client";

import { getPermalinkForThread } from "../utils";

export type Props = {
  thread: ThreadReference;
  editingEnabled?: boolean;
  movingEnabled?: boolean;
};

export function useThreadMenu({
  thread,
  editingEnabled,
  movingEnabled,
}: Props) {
  const router = useRouter();
  const account = useSession();
  const [_, setEditing] = useQueryState("edit", parseAsBoolean);
  const [, copyToClipboard] = useCopyToClipboard();

  const { deleteThread, revalidate } = useFeedMutations();

  const {
    isConfirming: isConfirmingDelete,
    handleConfirmAction: handleConfirmDelete,
    handleCancelAction: handleCancelDelete,
  } = useConfirmation(handleDelete);

  const isSharingEnabled = useShare();
  const isEditingEnabled = canEditPost(thread, account) && editingEnabled;
  const isMovingEnabled = canEditPost(thread, account) && movingEnabled;
  const isDeletingEnabled =
    canDeletePost(thread, account) && thread.deletedAt === undefined;

  const permalink = getPermalinkForThread(thread.slug);

  async function handleCopyLink() {
    copyToClipboard(permalink);
  }

  async function handleShare() {
    await navigator.share({
      title: `A post by ${thread.author.name}`,
      url: permalink,
      text: thread.description,
    });
  }

  function handleEdit() {
    setEditing(true);
  }

  async function handleDelete() {
    await handle(
      async () => {
        await withUndo({
          message: "Thread deleted",
          duration: 5000,
          toastId: `thread-${thread.id}`,
          action: async () => {
            await deleteThread(thread.id);
            router.push("/");
          },
          onUndo: () => {},
        });
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  return {
    isSharingEnabled,
    isEditingEnabled,
    isMovingEnabled,
    isDeletingEnabled,
    isConfirmingDelete,
    handlers: {
      handleCopyLink,
      handleShare,
      handleEdit,
      handleConfirmDelete,
      handleCancelDelete,
    },
  };
}
