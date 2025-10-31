"use client";

import { Reply, Thread } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { useShare } from "src/utils/client";

import { handle } from "@/api/client";
import { useThreadMutations } from "@/lib/thread/mutation";
import { withUndo } from "@/lib/thread/undo";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

import { getPermalinkForPost } from "../utils";

export type Props = {
  thread: Thread;
  reply: Reply;
  onEdit: () => void;
};

export function useReplyMenu({ thread, reply, onEdit }: Props) {
  const { revalidate, deleteReply } = useThreadMutations(thread);

  const account = useSession();
  const [, copyToClipboard] = useCopyToClipboard();

  const permalink = getPermalinkForPost(thread.slug, reply.id);

  const isSharingEnabled = useShare();
  const isEditingEnabled = account?.id === reply.author.id;
  const isDeletingEnabled = account?.id === reply.author.id;

  async function handleCopyURL() {
    copyToClipboard(permalink);
  }

  async function handleShare() {
    await navigator.share({
      title: `A post by ${reply.author.name}`,
      url: permalink,
      text: reply.body,
    });
  }

  function handleSetEditing() {
    onEdit();
  }

  async function handleDelete() {
    await handle(
      async () => {
        await withUndo({
          message: "Message deleted",
          duration: 5000,
          toastId: `reply-${reply.id}`,
          action: async () => {
            await deleteReply(reply.id);
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
    isDeletingEnabled,
    handlers: {
      handleCopyURL,
      handleShare,
      handleSetEditing,
      handleDelete,
    },
  };
}
