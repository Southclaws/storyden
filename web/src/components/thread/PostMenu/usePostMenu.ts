"use client";

import { useCopyToClipboard } from "@uidotdev/usehooks";
import { useRouter } from "next/navigation";
import { mutate } from "swr";

import { postDelete } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { getThreadGetKey, threadDelete } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { isShareEnabled } from "src/utils/client";
import { useToast } from "src/utils/useToast";

import { useThreadScreenContext } from "../context/context";
import { getPermalinkForPost } from "../utils";

export function usePostMenu(props: PostProps) {
  const router = useRouter();
  const toast = useToast();
  const account = useSession();
  const { thread, setEditingPostID } = useThreadScreenContext();
  const [, copyToClipboard] = useCopyToClipboard();

  const permalink = getPermalinkForPost(props.root_slug, props.id);

  const shareEnabled = isShareEnabled();
  const editEnabled = account?.id === props.author.id;
  const deleteEnabled = account?.id === props.author.id;

  async function onCopyLink() {
    copyToClipboard(permalink);
  }

  async function onShare() {
    await navigator.share({
      title: `A post by ${props.author.name}`,
      url: `#${props.id}`,
      text: props.body,
    });
  }

  async function onEdit() {
    setEditingPostID(props.id);
  }

  async function onDelete() {
    if (props.id === thread?.id) {
      await threadDelete(thread.id);
      toast({ title: "Thread deleted" });
      router.push("/");
    } else {
      await postDelete(props.id);
      toast({ title: "Post deleted" });
      mutate(getThreadGetKey(thread?.slug ?? props.id));
    }
  }

  return {
    onCopyLink,
    shareEnabled,
    onShare,
    editEnabled,
    onEdit,
    deleteEnabled,
    onDelete,
  };
}
