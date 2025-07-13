"use client";

import { useThreadList } from "@/api/openapi-client/threads";
import { Account, Category, ThreadListResult } from "@/api/openapi-schema";

export type Props = {
  initialPage?: number;
  initialPageData?: ThreadListResult;
  category?: Category;
};

export function useThreadFeedScreen(props: Props) {
  const initialPage = props.initialPage ?? 1;
  const { data, error } = useThreadList(
    {
      page: initialPage.toString(),
      categories: props.category ? [props.category.slug] : [],
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

  const showPaginationTop = data?.next_page && initialPage > 1;

  return {
    ready: true as const,
    showPaginationTop,
    data,
  };
}
