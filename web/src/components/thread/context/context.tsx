"use client";

import { createContext, useContext } from "react";

import { Identifier, Thread } from "src/api/openapi-schema";

export type ThreadScreenContextShape = {
  thread: Thread | undefined;

  editingPostID: Identifier | undefined;
  setEditingPostID: (postID: Identifier | undefined) => void;

  editingTitle: string | undefined;
  setEditingTitle: (title: string) => void;

  editingContent: string | undefined;
  setEditingContent: (content: string) => void;
};

const noop = (_: unknown) => {
  _;
};

export const ThreadScreenContext = createContext<ThreadScreenContextShape>({
  thread: undefined,
  editingPostID: undefined,
  setEditingPostID: noop,
  editingTitle: undefined,
  setEditingTitle: noop,
  editingContent: undefined,
  setEditingContent: noop,
});

export function useThreadScreenContext() {
  return useContext(ThreadScreenContext);
}
