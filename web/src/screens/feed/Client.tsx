"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";
import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { Unready } from "src/components/site/Unready";

export function Client(props: { category: string; threads: ThreadList }) {
  const { data, error } = useThreadList(
    { categories: [props.category] },
    {
      swr: {
        fallbackData: { threads: props.threads },
      },
    },
  );

  if (!data) return <Unready {...error} />;

  return <MixedPostList posts={data?.threads} />;
}
