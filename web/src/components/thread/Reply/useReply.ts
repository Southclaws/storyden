import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Reply, Thread } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useThreadMutations } from "@/lib/thread/mutation";

export const FormSchema = z.object({
  body: z.string().min(1, "Reply is empty."),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  thread: Thread;
  reply: Reply;
};

export function useReply({ thread, reply }: Props) {
  const { revalidate, updateReply } = useThreadMutations(thread);
  const [resetKey, setResetKey] = useState("");
  const [isEditing, setEditing] = useState(false);
  const [isEmpty, setEmpty] = useState(true);
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      body: reply.body,
    },
  });

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
    setResetKey(Date.now().toString());
  }

  const handleSave = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateReply(reply.id, data);

        setEditing(false);
        form.reset(data);
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Saved!",
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  });

  return {
    isEmpty,
    isEditing,
    resetKey,
    form,
    handlers: {
      handleSetEditing,
      handleEmptyStateChange,
      handleDiscardChanges,
      handleSave,
    },
  };
}
