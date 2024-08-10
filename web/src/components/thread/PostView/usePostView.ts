import { postUpdate } from "src/api/openapi-client/posts";
import { threadUpdate, useThreadGet } from "src/api/openapi-client/threads";
import { Post } from "src/api/openapi-schema";

import { useThreadScreenContext } from "../context/context";

export function usePostView(props: Post) {
  const {
    thread,
    editingPostID,
    setEditingPostID,
    editingTitle,
    editingContent,
    setEditingContent,
  } = useThreadScreenContext();

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
      });
    } else {
      await postUpdate(editingPostID, {
        body: editingContent,
      });
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
