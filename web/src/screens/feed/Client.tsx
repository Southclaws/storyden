"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";
import { Unready } from "src/components/Unready";
import { TextPostList } from "src/components/feed/text/TextPostList";

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

  return <TextPostList posts={data?.threads} />;
}
