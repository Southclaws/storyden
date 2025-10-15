"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { DatagraphItemKind, Thread } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { sendBeacon } from "@/lib/beacon/beacon";
import { useThreadMutations } from "@/lib/thread/mutation";
import { scrollToBottom } from "@/utils/scroll";

type Value = {
  body: string;
  isEmpty: boolean;
};

export const FormSchema = z.object({
  body: z.string().min(1, "Please enter a message."),
});
export type Form = z.infer<typeof FormSchema>;

export function useReplyBox(thread: Thread) {
  const session = useSession();
  const { createReply, revalidate } = useThreadMutations(thread);
  const [resetKey, setResetKey] = useState("");
  const [isEmpty, setEmpty] = useState(true);
  const form = useForm<Form>({ resolver: zodResolver(FormSchema) });

  function handleEmptyStateChange(isEmpty: boolean) {
    setEmpty(isEmpty);
  }

  const handleSubmit = form.handleSubmit(async (data: Form) => {
    await handle(
      async () => {
        await createReply(data);

        // Mark the thread as read after successfully replying to it
        try {
          sendBeacon(DatagraphItemKind.thread, thread.id);
        } catch (error) {
          console.warn("failed to send beacon:", error);
        }

        // This is a little hack tbh, essentially if this prop for the
        // ContentComposer component changes, its value is reset. Could have
        // done it with a hook but... meh this is simpler (albeit not idiomatic)
        setResetKey(new Date().toISOString());
        form.reset();
        setEmpty(true);

        scrollToBottom();
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  });

  return {
    isLoggedIn: !!session,
    isEmpty,
    isLoading: form.formState.isSubmitting,
    resetKey,
    form,
    handlers: {
      handleSubmit,
      handleEmptyStateChange,
    },
  };
}
