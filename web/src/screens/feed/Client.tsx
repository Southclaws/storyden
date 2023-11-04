"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { useFeed } from "src/components/feed/useFeed";
import { Unready } from "src/components/site/Unready";

export function Client(props: { category: string; threads: ThreadList }) {
  const { data, error, handlers } = useFeed(
    {
      categories: [props.category],
    },
    props.threads,
  );

  if (!data) return <Unready {...error} />;

  return (
    <MixedPostList posts={data?.threads} onDelete={handlers.handleDelete} />
  );
}
