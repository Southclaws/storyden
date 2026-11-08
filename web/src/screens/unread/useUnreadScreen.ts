import { parseAsInteger, useQueryState } from "nuqs";

import { useThreadList } from "@/api/openapi-client/threads";
import type { Account, ThreadListResult } from "@/api/openapi-schema";

export type Props = {
  initialSession: Account;
  initialPage: number;
  initialThreadList: ThreadListResult;
};

export function useUnreadScreen({ initialPage, initialThreadList }: Props) {
  const [page, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: initialPage,
  });

  const { data, error } = useThreadList(
    { page: page.toString(), read_status: "unread", ignore_pinned: true },
    {
      swr: {
        fallbackData: page === initialPage ? initialThreadList : undefined,
        dedupingInterval: 0,
        onSuccess(data) {
          const lastPage = Math.max(1, data.total_pages);
          if (page > lastPage) {
            void setPage(lastPage);
          }
        },
      },
    },
  );

  if (!data) {
    return { ready: false as const, error };
  }

  function handlePageChange(page: number) {
    void setPage(page);
  }

  return { ready: true as const, data, handlePageChange };
}
