"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { useFeed } from "src/components/feed/useFeed";
import { Onboarding } from "src/components/site/Onboarding/Onboarding";
import { Unready } from "src/components/site/Unready";

type Props = { threads: ThreadList };

export function Client(props: Props) {
  const { data, error, handlers } = useFeed(undefined, props.threads);

  if (!data) return <Unready {...error} />;

  return (
    <>
      <Onboarding />
      <MixedPostList posts={data?.threads} onDelete={handlers.handleDelete} />
    </>
  );
}
