"use client";

import { MixedContentFeed } from "src/components/feed/mixed/MixedContentFeed";
import { Props, useFeed } from "src/components/feed/useFeed";
import { Unready } from "src/components/site/Unready";

export function FeedScreenClient(props: Props) {
  const { data, error } = useFeed(props);

  if (!data) return <Unready {...error} />;

  return <MixedContentFeed data={data} />;
}
