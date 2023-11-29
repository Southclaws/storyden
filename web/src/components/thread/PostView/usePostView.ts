import { postUpdate } from "src/api/openapi/posts";
import { PostProps } from "src/api/openapi/schemas";
import { threadUpdate, useThreadGet } from "src/api/openapi/threads";
import { useToast } from "src/theme/components";

import { useThreadScreenContext } from "../context/context";

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
  const { mutate } = useThreadGet(thread!.slug);

  const isEditing = editingPostID === props.id;
  const isEditingThread = thread?.id === editingPostID;

  function onContentChange(value: string) {
    setEditingContent(value);
  }

  async function onPublishEdit() {
    if (!editingPostID) {
      throw new Error(
        "Cannot publish edits as the editing context has lost the target post ID.",
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

    await mutate();
    setEditingPostID(undefined);
  }

  function onCancelEdit() {
    setEditingPostID(undefined);
  }

  return {
    isEditing,
    editingContent,
    onContentChange,
    onPublishEdit,
    onCancelEdit,
  };
}
