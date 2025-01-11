"use client";

import { useThreadList } from "@/api/openapi-client/threads";
import { Account, Category, ThreadListResult } from "@/api/openapi-schema";

export type Props = {
  initialSession?: Account;
  initialPage: number;
  initialPageData?: ThreadListResult;
  category?: Category;
};

export function useThreadFeedScreen(props: Props) {
  const { data, error } = useThreadList(
    {
      page: props.initialPage.toString(),
    },
    {
      swr: {
        fallbackData: props.initialPageData,
      },
    },
  );
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const showPaginationTop = data?.next_page && props.initialPage > 1;

  return {
    ready: true as const,
    showPaginationTop,
    data,
  };
}
