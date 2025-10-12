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
  hideCategoryBadge = false,
  showQuickShare = true,
}: Props & {
  showCategorySelect: boolean;
  hideCategoryBadge?: boolean;
  showQuickShare?: boolean;
}) {
  const { session } = useSettingsContext();

  return (
    <LStack>
      {showQuickShare && (
        <QuickShare
          initialSession={session}
          initialCategory={category}
          showCategorySelect={showCategorySelect}
        />
      )}
      <ThreadFeed
        initialPage={initialPage}
        initialPageData={initialPageData}
        category={category}
        paginationBasePath={paginationBasePath}
        hideCategoryBadge={hideCategoryBadge}
      />
    </LStack>
  );
}

export function ThreadFeed(props: Props & { hideCategoryBadge?: boolean }) {
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
          return <ThreadReferenceCard key={t.slug} thread={t} hideCategoryBadge={props.hideCategoryBadge} />;
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
