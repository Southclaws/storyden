"use client";

import { Unready } from "src/components/Unready";

import { ThreadList } from "./components/ThreadList";
import { useContent } from "./useContent";

export function Content() {
  const { data, error } = useContent();

  if (!data) return <Unready {...error} />;

  return <ThreadList threads={data.threads} />;
}
