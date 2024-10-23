import { zodResolver } from "@hookform/resolvers/zod";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { threadUpdate, useThreadGet } from "@/api/openapi-client/threads";
import { ThreadGetResponse } from "@/api/openapi-schema";
import { useThreadSubscription } from "@/lib/thread/subscription";

export type Props = {
  slug: string;
  thread: ThreadGetResponse;
};

export const FormSchema = z.object({
  title: z.string().min(1, "Please enter a title."),
  body: z.string().min(1),
});
export type Form = z.infer<typeof FormSchema>;

export function useThreadScreen({ slug, thread }: Props) {
  useThreadSubscription(thread);

  const [editing, setEditing] = useQueryState("edit", parseAsBoolean);
  const [resetKey, setResetKey] = useState("");
  const [isEmpty, setEmpty] = useState(true);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    reValidateMode: "onChange",
    defaultValues: {
      title: thread.title,
      body: thread.body,
    },
  });

  const { data, error, mutate } = useThreadGet(slug, {
    swr: {
      fallbackData: thread,
    },
  });

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
    form.reset(thread);
    setEditing(false);
    setResetKey(Date.now().toString());
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

  return {
    ready: true as const,
    isEditing: editing,
    isEmpty,
    resetKey,
    form,
    data: {
      thread: data,
    },
    handlers: {
      handleEditing,
      handleEmptyStateChange,
      handleDiscardChanges,
      handleSave,
    },
  };
}
