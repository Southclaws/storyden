"use client";

import { useThreadGet } from "src/api/openapi-client/threads";
import { Thread } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";
import { ThreadView } from "src/components/thread/ThreadView/ThreadView";

export function Client(props: { slug: string; thread: Thread }) {
  const { data, error } = useThreadGet(props.slug, {
    swr: {
      fallbackData: props.thread,
    },
  });

  if (!data) return <Unready error={error} />;

  return <ThreadView {...data} />;
}
