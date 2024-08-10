"use client";

import { useState } from "react";

import { Thread } from "src/api/openapi/schemas";

import { ThreadScreenContextShape } from "./context";

export function useThreadScreenState(
  props: Thread | undefined,
): ThreadScreenContextShape {
  const [editingPostID, _setEditingPostID] =
    useState<ThreadScreenContextShape["editingPostID"]>(undefined);

  const [editingTitle, setEditingTitle] =
    useState<ThreadScreenContextShape["editingTitle"]>(undefined);

  const [editingContent, setEditingContent] =
    useState<ThreadScreenContextShape["editingContent"]>(undefined);

  function setEditingPostID(id: string | undefined) {
    _setEditingPostID(id);

    if (id === props?.id) {
      setEditingTitle(props?.title);
    }

    const target = props?.replies.find((p) => p.id === id);
    if (target) {
      setEditingContent(target.body);
    }
  }

  return {
    thread: props,
    editingPostID,
    setEditingPostID,
    editingTitle,
    setEditingTitle,
    editingContent,
    setEditingContent,
  };
}
