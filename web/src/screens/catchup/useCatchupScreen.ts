import {
  parseAsInteger,
  parseAsString,
  parseAsStringLiteral,
  useQueryState,
} from "nuqs";
import { useEffect, useRef } from "react";

import { useThreadList } from "@/api/openapi-client/threads";
import type {
  Account,
  CategoryListOKResponse,
  ThreadListResult,
} from "@/api/openapi-schema";

import {
  type CatchupView,
  catchupViews,
  getCatchupThreadParams,
} from "./query";

export type Props = {
  initialSession: Account;
  initialPage: number;
  initialView: CatchupView;
  initialCategory: string;
  initialThreadList: ThreadListResult;
  initialCategoryList: CategoryListOKResponse;
};

export function useCatchupScreen(props: Props) {
  const hydrated = useRef(false);
  const [page, setPage] = useQueryState(
    "page",
    parseAsInteger.withDefault(1).withOptions({ history: "push" }),
  );
  const [view, setView] = useQueryState(
    "view",
    parseAsStringLiteral(catchupViews)
      .withDefault("unread")
      .withOptions({ history: "push" }),
  );
  const [category, setCategory] = useQueryState(
    "category",
    parseAsString.withDefault("").withOptions({ history: "push" }),
  );

  const isInitialQuery =
    page === props.initialPage &&
    view === props.initialView &&
    category === props.initialCategory;

  const { data, error, mutate, isValidating } = useThreadList(
    getCatchupThreadParams(page, view, category),
    {
      swr: {
        fallbackData: isInitialQuery ? props.initialThreadList : undefined,
        revalidateOnMount: !isInitialQuery,
        revalidateIfStale: true,
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
        dedupingInterval: 0,
        onSuccess(data) {
          const lastPage = Math.max(1, data.total_pages);
          if (page > lastPage) {
            void setPage(lastPage, { history: "replace" });
          }
        },
      },
    },
  );

  // Populate the shared cache so thread actions can patch the SSR list
  // without refetching it and changing its membership.
  useEffect(() => {
    if (isInitialQuery && !hydrated.current) {
      hydrated.current = true;
      void mutate((cached) => cached ?? props.initialThreadList, {
        revalidate: false,
      });
    }
  }, [isInitialQuery, props.initialThreadList, mutate]);

  function handlePageChange(page: number) {
    void setPage(page);
  }

  function handleViewChange(view: string) {
    if (view === "unread" || view === "unresolved") {
      void setView(view);
      void setPage(1);
    }
  }

  function handleCategoryChange(category: string) {
    void setCategory(category);
    void setPage(1);
  }

  return {
    data,
    error,
    view,
    category,
    isValidating,
    refresh: mutate,
    handlePageChange,
    handleViewChange,
    handleCategoryChange,
  };
}
