"use client";

import { ThreadList } from "src/api/openapi/schemas";
import { useThreadList } from "src/api/openapi/threads";
import { TextPostList } from "src/components/feed/text/TextPostList";
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
      <TextPostList posts={data?.threads} />
    </>
  );
}
