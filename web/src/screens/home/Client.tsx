"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";
import { MixedPostList } from "src/components/feed/mixed/MixedPostList";
import { Onboarding } from "src/components/site/Onboarding/Onboarding";
import { Unready } from "src/components/site/Unready";

type Props = { threads: ThreadList };

export function Client(props: Props) {
  const { data, error } = useThreadList(
    {},
    {
      swr: {
        fallbackData: { threads: props.threads },
      },
    },
  );

  if (!data) return <Unready {...error} />;

  return (
    <>
      <Onboarding />
      <MixedPostList posts={data?.threads} />
    </>
  );
}
