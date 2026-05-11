import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Account, Permission, Reply, Thread } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useConfirmation } from "@/components/site/useConfirmation";
import { useI18n } from "@/i18n/provider";
import { useReportContext } from "@/lib/report/useReportContext";
import type { SignatureConfig } from "@/lib/settings/settings";
import { useThreadMutations } from "@/lib/thread/mutation";
import { withUndo } from "@/lib/thread/undo";
import { hasPermission } from "@/utils/permissions";

function getFormSchema(t: (key: string) => string) {
  return z.object({
    body: z.string().min(1, t("Reply is empty.")),
  });
}
export type Form = z.infer<ReturnType<typeof getFormSchema>>;

export type Props = {
  initialSession?: Account;
  thread: Thread;
  reply: Reply;
  currentPage?: number;
  initialSignatureConfig: SignatureConfig;
};

export function useReply({
  initialSession,
  thread,
  reply,
  currentPage,
}: Props) {
  const { t } = useI18n();
  const session = useSession(initialSession);
  const { resolveReport } = useReportContext();
  const { revalidate, updateReply, deleteReply } = useThreadMutations(
    thread,
    currentPage,
    undefined,
  );
  const [resetKey, setResetKey] = useState("");
  const [isEditing, setEditing] = useState(false);
  const [isEditingInReview, setEditingInReview] = useState(false);
  const [isEmpty, setEmpty] = useState(true);
  const form = useForm<Form>({
    resolver: zodResolver(getFormSchema(t)),
    defaultValues: {
      body: reply.body,
    },
  });

  const {
    isConfirming: isConfirmingDelete,
    handleConfirmAction: handleConfirmDelete,
    handleCancelAction: handleCancelDelete,
  } = useConfirmation(handleDeleteReply);

  const canManageReplies = hasPermission(session, Permission.MANAGE_POSTS);

  function handleEmptyStateChange(isEmpty: boolean) {
    setEmpty(isEmpty);
  }

  function handleSetEditing() {
    setEditing(true);
  }

  function handleDiscardChanges() {
    // TODO: useConfirmation
    form.reset(reply);
    setEditing(false);
    setEditingInReview(false);
    setResetKey(Date.now().toString());
  }

  function handleSetEditingInReview() {
    setEditingInReview(true);
    setEditing(true);
  }

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const updates = isEditingInReview
          ? { ...data, visibility: "published" as const }
          : data;

        await updateReply(reply.id, updates);

        if (isEditingInReview) {
          await resolveReport();
        }

        setEditing(false);
        setEditingInReview(false);
        form.reset(data);
      },
      {
        promiseToast: {
          loading: t("Saving..."),
          success: isEditingInReview ? t("Saved and published!") : t("Saved!"),
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  });

  async function handleAcceptReply() {
    await handle(
      async () => {
        await updateReply(reply.id, { visibility: "published" });
        await resolveReport();
      },
      {
        promiseToast: {
          loading: t("Accepting..."),
          success: t("Reply accepted!"),
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  }

  async function handleDeleteReply() {
    await handle(
      async () => {
        await withUndo({
          message: t("Reply deleted"),
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
    isEmpty,
    isEditing,
    isEditingInReview,
    canManageReplies,
    resetKey,
    form,
    isConfirmingDelete,
    handlers: {
      handleSetEditing,
      handleSetEditingInReview,
      handleEmptyStateChange,
      handleDiscardChanges,
      handleSave,
      handleAcceptReply,
      handleConfirmDelete,
      handleCancelDelete,
    },
  };
}
