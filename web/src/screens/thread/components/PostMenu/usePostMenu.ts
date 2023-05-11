import { PostProps } from "src/api/openapi/schemas";
import { useThreadScreenContext } from "../../context";
import { useAuthProvider } from "src/auth/useAuthProvider";
import { useClipboard } from "@chakra-ui/react";
import { getPermalinkForPost } from "../../utils";

export function usePostMenu(props: PostProps) {
  const { account } = useAuthProvider();
  const { setEditingPostID } = useThreadScreenContext();
  const { onCopy } = useClipboard(
    getPermalinkForPost(props.root_slug, props.id)
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
    console.log("byebye");
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
