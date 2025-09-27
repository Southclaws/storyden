"use client";

import { parseAsInteger, useQueryState } from "nuqs";

import { useThreadList } from "@/api/openapi-client/threads";
import { Category, ThreadListResult } from "@/api/openapi-schema";

export type Props = {
  initialPage?: number;
  initialPageData?: ThreadListResult;
  category:
    | undefined // No category specified, no filters applied.
    | Category // An explicit category.
    | null; // Explicitly uncategorised.
  paginationBasePath: string;
};

export function useThreadFeedScreen(props: Props) {
  const initialPage = props.initialPage ?? 1;

  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: props.initialPage ?? 1,
  });

  function handlePageChange(page: number) {
    setPage(page);
  }

  const { data, error } = useThreadList(
    {
      page: page.toString(),
      categories:
        props.category === undefined
          ? []
          : [props.category === null ? "null" : props.category.slug],
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
    handlePageChange,
  };
}
