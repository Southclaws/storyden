"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Account, DatagraphItemKind, Thread } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { sendBeacon } from "@/lib/beacon/beacon";
import { useThreadMutations } from "@/lib/thread/mutation";

import { useReplyContext } from "../ReplyContext";

export type Props = {
  initialSession?: Account;
  thread: Thread;
};

type ReplyLocationState = {
  id: string;
  pageNumber: number;
  permalink: string;
};

export const FormSchema = z.object({
  body: z.string().min(1, "Please enter a message."),
});
export type Form = z.infer<typeof FormSchema>;

export function useReplyBox({ initialSession, thread }: Props) {
  const session = useSession(initialSession);
  const { replyTo, clearReplyTo } = useReplyContext();
  const { createReply, revalidate } = useThreadMutations(
    thread,
    thread.replies.current_page,
    thread.replies.total_pages,
  );
  const [resetKey, setResetKey] = useState("");
  const [isEmpty, setEmpty] = useState(true);
  const [postedReply, setPostedReply] = useState<ReplyLocationState | null>(
    null,
  );
  const form = useForm<Form>({ resolver: zodResolver(FormSchema) });

  function handleEmptyStateChange(isEmpty: boolean) {
    setEmpty(isEmpty);
  }

  function handleReplyPostedAdmonitionClose() {
    setPostedReply(null);
  }

  function handleReplyNavigation() {
    setPostedReply(null);
  }

  const handleSubmit = form.handleSubmit(async (data: Form) => {
    await handle(
      async () => {
        const { id } = await createReply({
          body: data.body,
          reply_to: replyTo?.reply.id,
        });

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
        clearReplyTo();

        // If we are not on the last page, we need to inform the user that their
        // reply is on a different page and provide them a link to navigate.
        const currentPage = thread.replies.current_page;
        const totalPages = thread.replies.total_pages;
        const isLastPage =
          !currentPage || !totalPages || currentPage === totalPages;
        if (!isLastPage && totalPages) {
          setPostedReply({
            id,
            pageNumber: totalPages,
            permalink: `/t/${thread.slug}?page=${totalPages}#${id}`,
          });
        }
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
    postedReply,
    form,
    handlers: {
      handleSubmit,
      handleEmptyStateChange,
      handleReplyPostedAdmonitionClose,
      handleReplyNavigation,
    },
  };
}
