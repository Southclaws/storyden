"use client";

import { Unready } from "src/components/site/Unready";
import { ThreadIndexView } from "src/components/thread/ThreadIndexView/ThreadIndexView";

import { Props, useThreadIndexScreen } from "./useThreadIndexScreen";

export function ThreadIndexScreen(props: Props) {
  const { ready, data, mutate, error } = useThreadIndexScreen(props);

  if (!ready) return <Unready error={error} />;

  return (
    <ThreadIndexView
      threads={data}
      mutate={mutate}
      query={props.query}
      page={props.page}
    />
  );
}
