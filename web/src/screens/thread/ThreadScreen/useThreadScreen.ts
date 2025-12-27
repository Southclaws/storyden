import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { threadUpdate, useThreadGet } from "@/api/openapi-client/threads";
import {
  Account,
  DatagraphItemKind,
  Permission,
  ThreadGetResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useConfirmation } from "@/components/site/useConfirmation";
import { useBeacon } from "@/lib/beacon/useBeacon";
import { useReportContext } from "@/lib/report/useReportContext";
import { useThreadMutations } from "@/lib/thread/mutation";
import { withUndo } from "@/lib/thread/undo";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  initialSession?: Account;
  initialPage?: number;
  slug: string;
  thread: ThreadGetResponse;
};

export const FormSchema = z.object({
  title: z.string().min(1, "Please enter a title."),
  body: z.string().min(1),
  tags: z.array(z.string()).optional(),
});
export type Form = z.infer<typeof FormSchema>;

export function useThreadScreen({
  initialSession,
  initialPage,
  slug,
  thread,
}: Props) {
  const router = useRouter();
  const session = useSession(initialSession);
  const { reportId, resolveReport } = useReportContext();
  const { updateReply, deleteReply, revalidate } = useThreadMutations(
    thread,
    initialPage,
  );

  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });
  const [resetKey, setResetKey] = useState("");
  const [isEmpty, setEmpty] = useState(
    !thread.body || thread.body.trim().length === 0,
  );

  const {
    isConfirming: isConfirmingDelete,
    handleConfirmAction: handleConfirmDelete,
    handleCancelAction: handleCancelDelete,
  } = useConfirmation(handleDeleteThread);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    reValidateMode: "onChange",
    defaultValues: {
      title: thread.title,
      body: thread.body,
    },
  });

  const { data, error, mutate } = useThreadGet(
    slug,
    {
      page: initialPage?.toString(),
    },
    {
      swr: {
        fallbackData: thread,
      },
    },
  );

  useBeacon(DatagraphItemKind.thread, data?.id);

  const isModerator = hasPermission(
    session,
    Permission.MANAGE_POSTS,
    Permission.ADMINISTRATOR,
  );

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  function handleEditing() {
    setEditing(true);
  }

  function handleEmptyStateChange(isEmpty: boolean) {
    setEmpty(isEmpty);
  }

  function handleDiscardChanges() {
    // TODO: useConfirmation
    form.reset({
      title: thread.title,
      body: thread.body,
      tags: thread.tags.map((t) => t.name),
    });
    setEditing(false);
    setResetKey(Date.now().toString());
    setEmpty(!thread.body || thread.body.trim().length === 0);
  }

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await mutate({
          ...thread,
          title: data.title,
          body: data.body,
        });

        await threadUpdate(slug, data);

        setEditing(false);
        form.reset(data);
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Saved!",
        },
        cleanup: async () => {
          await mutate();
        },
      },
    );
  });

  async function handleAcceptThread() {
    await handle(
      async () => {
        await updateReply(thread.id, { visibility: "published" });
        await resolveReport();
      },
      {
        promiseToast: {
          loading: "Accepting...",
          success: "Thread accepted!",
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  }

  function handleEditAndAccept() {
    setEditing(true);
  }

  async function handleDeleteThread() {
    await handle(
      async () => {
        await withUndo({
          message: "Thread deleted",
          duration: 5000,
          toastId: `thread-${thread.id}`,
          action: async () => {
            await deleteReply(thread.id);
            await resolveReport();
            if (reportId) {
              router.push(`/reports`);
            } else {
              router.push("/");
            }
          },
          onUndo: () => {},
        });
      },
      {
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  }

  return {
    ready: true as const,
    isEditing: editing,
    isEmpty,
    resetKey,
    form,
    isModerator,
    isConfirmingDelete,
    data: {
      thread: data,
    },
    handlers: {
      handleEditing,
      handleEmptyStateChange,
      handleDiscardChanges,
      handleSave,
      handleAcceptThread,
      handleEditAndAccept,
      handleConfirmDelete,
      handleCancelDelete,
    },
  };
}
