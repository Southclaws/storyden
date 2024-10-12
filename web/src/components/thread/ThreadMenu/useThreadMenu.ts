"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";

import { Thread } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useFeedMutations } from "@/lib/feed/mutation";
import { useShare } from "@/utils/client";

import { getPermalinkForThread } from "../utils";

export type Props = {
  thread: Thread;
};

export function useThreadMenu({ thread }: Props) {
  const router = useRouter();
  const account = useSession();
  const [_, setEditing] = useQueryState("edit", parseAsBoolean);
  const [, copyToClipboard] = useCopyToClipboard();

  const { deleteThread, revalidate } = useFeedMutations();

  const isSharingEnabled = useShare();
  const isEditingEnabled = account?.id === thread.author.id;
  const isDeletingEnabled = account?.id === thread.author.id;

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

  function handleDelete() {
    handle(
      async () => {
        await deleteThread(thread.id);
        router.push("/");
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
      handleCopyLink,
      handleShare,
      handleEdit,
      handleDelete,
    },
  };
}
