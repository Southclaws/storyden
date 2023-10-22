import { useClipboard, useToast } from "@chakra-ui/react";
import { useRouter } from "next/navigation";
import { mutate } from "swr";

import { postDelete } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { getThreadGetKey, threadDelete } from "src/api/openapi/threads";
import { useSession } from "src/auth";

import { useThreadScreenContext } from "../../context";
import { getPermalinkForPost } from "../../utils";

export function usePostMenu(props: PostProps) {
  const router = useRouter();
  const toast = useToast();
  const account = useSession();
  const { thread, setEditingPostID } = useThreadScreenContext();
  const { onCopy } = useClipboard(
    getPermalinkForPost(props.root_slug, props.id),
  );

  const shareEnabled = !!navigator.share;
  const editEnabled = account?.id === props.author.id;
  const deleteEnabled = account?.id === props.author.id;

  async function onCopyLink() {
    onCopy();
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
