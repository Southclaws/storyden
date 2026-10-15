"use client";

import {
  type PropsWithChildren,
  type ReactNode,
  createContext,
  useContext,
} from "react";

const Context = createContext<Record<string, ReactNode>>({});

export function RenderedContentProvider({
  content,
  children,
}: PropsWithChildren<{ content: Record<string, ReactNode> }>) {
  const parent = useContext(Context);
  return <Context value={{ ...parent, ...content }}>{children}</Context>;
}

export function useRenderedContent(html: string) {
  const content = useContext(Context);
  return Object.hasOwn(content, html) ? content[html] : undefined;
}
