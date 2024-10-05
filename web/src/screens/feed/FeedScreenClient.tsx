"use client";

import { Props, useFeed } from "src/components/feed/useFeed";
import { Unready } from "src/components/site/Unready";

import { TextPostList } from "@/components/feed/text/TextPostList";

export function FeedScreenClient(props: Props) {
  const { data, error } = useFeed(props);

  if (!data) return <Unready error={error} />;

  return <TextPostList threads={data.threads} />;
}
