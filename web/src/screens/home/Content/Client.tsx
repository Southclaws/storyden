"use client";

import { Unready } from "src/components/Unready";

import { ThreadList } from "../components/ThreadList";

import { useClient } from "./useClient";

export function Client(props: { showEmptyState: boolean }) {
  const { data, error } = useClient();

  if (!data) return <Unready {...error} />;

  return (
    <>
      <ThreadList
        showEmptyState={props.showEmptyState}
        threads={data.threads}
      />
    </>
  );
}
