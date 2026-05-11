"use client";

import { Reply, Thread } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { useShare } from "src/utils/client";

import { handle } from "@/api/client";
import { useI18n } from "@/i18n/provider";
import { useReportContext } from "@/lib/report/useReportContext";
import { useThreadMutations } from "@/lib/thread/mutation";
import { canDeletePost, canEditPost } from "@/lib/thread/permissions";
import { withUndo } from "@/lib/thread/undo";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

import { getPermalinkForPost } from "../utils";

export type Props = {
  thread: Thread;
  reply: Reply;
  currentPage?: number;
  onEdit: () => void;
};

export function useReplyMenu({ thread, reply, currentPage, onEdit }: Props) {
  const { t } = useI18n();
  const { revalidate, deleteReply } = useThreadMutations(thread, currentPage);
  const { resolveReport } = useReportContext();

  const account = useSession();
  const [, copyToClipboard] = useCopyToClipboard();

  const permalink = getPermalinkForPost(thread.slug, reply.id, currentPage);

  const isSharingEnabled = useShare();
  const isEditingEnabled =
    canEditPost(reply, account) && reply.deletedAt === undefined;
  const isDeletingEnabled =
    canDeletePost(reply, account) && reply.deletedAt === undefined;

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
          message: t("Message deleted"),
          duration: 5000,
          toastId: `reply-${reply.id}`,
          action: async () => {
            await deleteReply(reply.id);
            await resolveReport();
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
