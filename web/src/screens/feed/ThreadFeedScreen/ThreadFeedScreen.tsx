"use client";

import { useIntersectionObserver } from "@uidotdev/usehooks";
import { parseAsInteger, useQueryState } from "nuqs";
import { PropsWithChildren, useEffect } from "react";

import { Unready } from "src/components/site/Unready";

import { FeedEmptyState } from "@/components/feed/FeedEmptyState";
import { ThreadReferenceCard } from "@/components/post/ThreadCard";
import { PaginationBubble } from "@/components/site/PaginationBubble/PaginationBubble";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import {
  Props,
  getPageElementAnchor,
  useThreadFeedScreen,
} from "./useThreadFeedScreen";

export function ThreadFeedScreen(props: Props) {
  const { ready, error, isValidating, morePagesAvailable, data, handlers } =
    useThreadFeedScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { pages, totalPages, pageSize } = data;

  if (pages.length === 0) {
    return <FeedEmptyState />;
  }

  return (
    <VStack w="full">
      <ol className={lstack()}>
        {pages.map((page) => {
          return (
            <PageSection
              key={getPageElementAnchor(page.current_page)}
              page={page.current_page}
            >
              {page.current_page > 1 && (
                <Heading color="fg.subtle" size="sm">
                  page {page.current_page}
                </Heading>
              )}

              <ol className={lstack()}>
                {page.threads.map((t) => {
                  return <ThreadReferenceCard key={t.slug} thread={t} />;
                })}
              </ol>
            </PageSection>
          );
        })}
      </ol>

      <VStack w="full">
        {isValidating ? (
          <Unready />
        ) : morePagesAvailable ? (
          <Button
            w="full"
            variant="outline"
            size="sm"
            onClick={handlers.handleNextPage}
          >
            Load more...
          </Button>
        ) : (
          <VStack
            w="full"
            color="fg.muted"
            textAlign="center"
            textWrap="balance"
          >
            <p>You&apos;ve reached the end.</p>
          </VStack>
        )}
      </VStack>

      <PaginationBubble
        path="/"
        totalPages={totalPages}
        pageSize={pageSize}
        onPageChange={handlers.handlePageChange}
      />
    </VStack>
  );
}

function PageSection({ children, page }: PropsWithChildren<{ page: number }>) {
  const [_, setPage] = useQueryState("page", {
    ...parseAsInteger,
    defaultValue: 1,
    clearOnDefault: true,
  });
  const [ref, entry] = useIntersectionObserver({
    threshold: 0.3,
    root: null,
    rootMargin: "0px",
  });

  useEffect(() => {
    if (entry?.isIntersecting) {
      setPage(page);
    }
  }, [setPage, entry, page]);

  return (
    <ol id={getPageElementAnchor(page)} ref={ref} className={lstack()}>
      {children}
    </ol>
  );
}
