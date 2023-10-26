"use client";

import { ChangeEvent } from "react";

import { Thread } from "src/api/openapi/schemas";

import { useThreadScreenContext } from "../context/context";

export function useTitle(thread: Thread) {
  const { editingPostID, editingTitle, setEditingTitle } =
    useThreadScreenContext();

  const editing = editingPostID === thread.id;

  function onTitleChange(e: ChangeEvent<HTMLInputElement>) {
    setEditingTitle(e.target.value);
  }

  return { editing, editingTitle, onTitleChange };
}
