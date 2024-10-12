"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";

import { Reply, Thread } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { useShare } from "src/utils/client";

import { handle } from "@/api/client";
import { useThreadMutations } from "@/lib/thread/mutation";

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

  function handleDelete() {
    handle(() => deleteReply(reply.id), {
      cleanup: async () => await revalidate(),
    });
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
