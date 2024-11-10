"use client";

import { Unready } from "src/components/site/Unready";

import { useThreadList } from "@/api/openapi-client/threads";
import { ThreadListParams, ThreadListResult } from "@/api/openapi-schema";
import { FeedEmptyState } from "@/components/feed/FeedEmptyState";
import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";

export type Props = {
  params?: ThreadListParams;
  initialData?: ThreadListResult;
};

export function ThreadFeedScreen({ params, initialData }: Props) {
  const { data, error } = useThreadList(params, {
    swr: { fallbackData: initialData },
  });
  if (!data) {
    return <Unready error={error} />;
  }

  if (data.threads.length === 0) {
    return <FeedEmptyState />;
  }

  return <ThreadReferenceList threads={data.threads} />;
}
