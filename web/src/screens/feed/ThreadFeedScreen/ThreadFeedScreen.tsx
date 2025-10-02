"use client";

import { Unready } from "src/components/site/Unready";

import { FeedEmptyState } from "@/components/feed/FeedEmptyState";
import { QuickShare } from "@/components/feed/QuickShare/QuickShare";
import { ThreadReferenceCard } from "@/components/post/ThreadCard";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { useSettingsContext } from "@/components/site/SettingsContext/SettingsContext";
import { LStack, VStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useThreadFeedScreen } from "./useThreadFeedScreen";

export function ThreadFeedScreen({
  initialPage,
  initialPageData,
  category,
  paginationBasePath,
  showCategorySelect,
}: Props & {
  showCategorySelect: boolean;
}) {
  const { session } = useSettingsContext();

  return (
    <LStack>
      <QuickShare
        initialSession={session}
        initialCategory={category}
        showCategorySelect={showCategorySelect}
      />
      <ThreadFeed
        initialPage={initialPage}
        initialPageData={initialPageData}
        category={category}
        paginationBasePath={paginationBasePath}
      />
    </LStack>
  );
}

export function ThreadFeed(props: Props) {
  const { ready, error, showPaginationTop, data, handlePageChange } =
    useThreadFeedScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  if (data.threads.length === 0) {
    return <FeedEmptyState />;
  }

  return (
    <VStack w="full">
      {showPaginationTop && (
        <PaginationControls
          path="/"
          currentPage={data.current_page}
          totalPages={data.total_pages}
          pageSize={data.page_size}
          onClick={handlePageChange}
        />
      )}
      <ol className={lstack()}>
        {data.threads.map((t) => {
          return <ThreadReferenceCard key={t.slug} thread={t} />;
        })}
      </ol>

      <PaginationControls
        path={props.paginationBasePath}
        currentPage={data.current_page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
        onClick={handlePageChange}
      />
    </VStack>
  );
}
