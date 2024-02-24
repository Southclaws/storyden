"use client";

import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { Props, useFeed } from "src/components/feed/useFeed";
import { Onboarding } from "src/components/site/Onboarding/Onboarding";
import { Unready } from "src/components/site/Unready";

export function FeedScreenClient(props: Props) {
  const { data, error, handlers } = useFeed(props);

  if (!data) return <Unready {...error} />;

  return (
    <>
      <Onboarding />
      <MixedPostList
        posts={data?.threads}
        onDelete={handlers.handleDeleteThread}
      />
    </>
  );
}
