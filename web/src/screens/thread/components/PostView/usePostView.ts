import { PostProps } from "src/api/openapi/schemas";
import { useThreadScreenContext } from "../../context";
import { threadUpdate } from "src/api/openapi/threads";
import { postUpdate } from "src/api/openapi/posts";
import { useToast } from "@chakra-ui/react";

export function usePostView(props: PostProps) {
  const {
    thread,
    editingPostID,
    setEditingPostID,
    editingTitle,
    editingContent,
    setEditingContent,
  } = useThreadScreenContext();
  const toast = useToast();

  const isEditing = editingPostID === props.id;
  const isEditingThread = thread?.id === editingPostID;

  async function onPublishEdit() {
    if (!editingPostID) {
      throw new Error(
        "Cannot publish edits as the editing context has lost the target post ID."
      );
    }

    if (isEditingThread) {
      await threadUpdate(editingPostID, {
        title: editingTitle,
        body: editingContent,
      }).then(() => toast({ title: "Thread updated" }));
    } else {
      await postUpdate(editingPostID, {
        body: editingContent,
      }).then(() => toast({ title: "Post updated" }));
    }

    setEditingPostID(undefined);
  }

  function onCancelEdit() {
    setEditingPostID(undefined);
  }

  return {
    isEditing,
    editingContent,
    setEditingContent,
    onPublishEdit,
    onCancelEdit,
  };
}
