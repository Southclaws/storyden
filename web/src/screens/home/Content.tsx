"use client";

import { Unready } from "src/components/Unready";

import { ThreadList } from "./components/ThreadList";
import { useContent } from "./useContent";

export function Content(props: { showEmptyState: boolean }) {
  const { data, error } = useContent();

  if (!data) return <Unready {...error} />;

  return (
    <ThreadList showEmptyState={props.showEmptyState} threads={data.threads} />
  );
}
