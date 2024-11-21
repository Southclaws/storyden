"use client";

import { last } from "lodash";
import { useCallback, useEffect, useState } from "react";
import useSWRInfinite from "swr/infinite";

import { fetcher } from "@/api/client";
import {
  APIError,
  Category,
  ThreadListOKResponse,
  ThreadListResult,
} from "@/api/openapi-schema";
import { getThreadListPageKey } from "@/lib/feed/mutation";

export type Props = {
  initialPage: number;
  initialPageData?: ThreadListResult[];
  category?: Category;
};

export function getPageElementAnchor(page: number) {
  return `page_${page}`;
}

// NOTE: Orval's infinite pagination is broken so this is done manually.
const threadListFetcher = (url: string) =>
  fetcher<ThreadListOKResponse>({ url });

export function useThreadFeedScreen(props: Props) {
  const parameters: Record<string, any> = {
    ...(props.category ? { categories: [props.category.slug] } : {}),
  };

  const loadingIntoFurtherPage = props.initialPage > 1;

  const key = getThreadListPageKey(parameters);

  // waitingToScrollTo is used to defer scrolling to a page until it has loaded
  // because if you've only loaded page 1 and 2 but want to scroll to page 5,
  // we first run setSize to trigger useSWRInfinite to start requesting those
  // pages (3, 4 and 5) and set this state to 5. Once the pages have loaded the
  // useEffect hook below triggers the scroll to the new rendered page section.
  const [waitingToScrollTo, setWaitingToScrollTo] = useState<number | null>(
    loadingIntoFurtherPage ? props.initialPage : null,
  );

  const { data, error, isLoading, isValidating, mutate, size, setSize } =
    useSWRInfinite<ThreadListOKResponse, APIError>(key, threadListFetcher, {
      keepPreviousData: false,
      initialSize: props.initialPage,
      fallbackData: props.initialPageData,
    });

  // This useEffect deals purely with deferred scroll-to behaviour, this is
  // triggered when waitingToScrollTo is set to a page number while the swr hook
  // is loading the next n pages in order to get to this one. It's not currently
  // very efficient as it requires loading all n pages before being able to get
  // to the one the user wants. In future we can use virtual scrolling to fix.
  useEffect(() => {
    if (waitingToScrollTo == null) {
      // not performing an awaited scroll-to
      return;
    }

    if (isLoading || isValidating) {
      // still loading the next page...
      return;
    }

    if (waitingToScrollTo > size) {
      // the swr size is not big enough to hold the page we want to scroll to.
      return;
    }

    const element = document.getElementById(
      getPageElementAnchor(waitingToScrollTo),
    );

    if (!element) {
      // Page section has gone missing somehow
      console.warn("unable to find page section for newly loaded scrollTo", {
        waitingToScrollTo,
        size,
        id: getPageElementAnchor(size),
      });
      return;
    }

    element?.scrollIntoView({
      behavior: "instant",
    });

    setWaitingToScrollTo(null);
  }, [isLoading, isValidating, size, waitingToScrollTo]);

  const morePagesAvailable = Boolean(last(data)?.next_page);

  // handleNextPage bumps useSWRInfinite to the next page
  const handleNextPage = useCallback(() => {
    if (!morePagesAvailable) {
      return;
    }

    setSize(size + 1);
  }, [morePagesAvailable, setSize, size]);

  // handleScroll deals with the infinite scroll behavior.
  const handleScroll = useCallback(() => {
    if (isLoading) {
      // Don't attempt a scroll-triggered page while already loading
      return;
    }

    if (!morePagesAvailable) {
      // Don't attempt a scroll-triggered page if there are no more pages
      return;
    }

    const scrolledToEnd =
      window.innerHeight + document.documentElement.scrollTop ===
      document.documentElement.offsetHeight;

    if (!scrolledToEnd) {
      return;
    }

    handleNextPage();
  }, [handleNextPage, morePagesAvailable, isLoading]);

  const handlePageChange = useCallback(
    (page: number) => {
      const element = document.getElementById(getPageElementAnchor(page));

      if (element) {
        // This page has been loaded already
        element?.scrollIntoView({
          behavior: "instant",
        });
      } else {
        window.scroll({
          behavior: "instant",
          top: document.documentElement.scrollHeight,
        });
        // This page has not been loaded yet
        // - Set the page to trigger useSWRInfinite to load it
        // - wait for the content to load
        // - scroll to the page section
        setSize(page);
        setWaitingToScrollTo(page);
      }
    },
    [setSize],
  );

  useEffect(() => {
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [handleScroll, isLoading, morePagesAvailable]);

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const lastPage = last(data)!;

  return {
    ready: true as const,
    isValidating,
    morePagesAvailable,
    data: {
      pages: data,
      totalPages: lastPage.total_pages,
      pageSize: lastPage.page_size,
    },
    handlers: {
      handleNextPage,
      handlePageChange,
    },
  };
}
