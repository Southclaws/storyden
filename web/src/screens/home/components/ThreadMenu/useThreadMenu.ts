"use client";

import { useClipboard, useToast } from "@chakra-ui/react";
import { mutate } from "swr";

import { ThreadReference } from "src/api/openapi/schemas";
import { getThreadListKey, threadDelete } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import {
  getPermalinkForThread,
  useQueryParameters,
} from "src/screens/home/utils";

export function useThreadMenu(props: ThreadReference) {
  const toast = useToast();
  const account = useSession();
  const { onCopy } = useClipboard(getPermalinkForThread(props.slug));
  const { category } = useQueryParameters();

  const shareEnabled = !!navigator.share;
  const deleteEnabled = account?.id === props.author.id;

  async function onCopyLink() {
    onCopy();
  }

  async function onShare() {
    await navigator.share({
      title: `A post by ${props.author.name}`,
      url: `#${props.id}`,
      text: props.short,
    });
  }

  async function onDelete() {
    await threadDelete(props.id);
    toast({ title: "Thread deleted" });
    mutate(
      getThreadListKey({
        categories: category ? [category] : undefined,
      })
    );
  }

  return {
    onCopyLink,
    shareEnabled,
    onShare,
    deleteEnabled,
    onDelete,
  };
}
