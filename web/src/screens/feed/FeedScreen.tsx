"use client";

import { Unready } from "src/components/site/Unready";

import { useThreadList } from "@/api/openapi-client/threads";
import { ThreadListParams, ThreadListResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/feed/EmptyState";
import { ThreadItemList } from "@/components/feed/ThreadItemList";

export type Props = {
  params?: ThreadListParams;
  initialData?: ThreadListResult;
};

export function FeedScreen({ params, initialData }: Props) {
  const { data, error } = useThreadList(params, {
    swr: { fallbackData: initialData },
  });
  if (!data) {
    return <Unready error={error} />;
  }

  if (data.threads.length === 0) {
    return <EmptyState />;
  }

  return <ThreadItemList threads={data.threads} />;
}
