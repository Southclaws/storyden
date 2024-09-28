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

export const ThreadScreenContext =
  createContext<ThreadScreenContextShape | null>(null);

export function useThreadScreenContext() {
  const ctx = useContext(ThreadScreenContext);

  if (!ctx) {
    throw new Error(
      "useThreadScreenContext must be used within a ThreadScreenContext.Provider",
    );
  }

  return ctx;
}
