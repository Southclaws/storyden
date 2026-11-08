"use client";

import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { CheckIcon } from "@/components/ui/icons/Check";
import { PageHeader } from "@/components/ui/page-header";
import { LStack } from "@/styled-system/jsx";

import { type Props, useUnreadScreen } from "./useUnreadScreen";

export function UnreadScreen(props: Props) {
  const { ready, data, error, handlePageChange } = useUnreadScreen(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <LStack className="unread-page" data-sd-page="unread" gap="4">
      <PageHeader title="Unread threads" />

      {data.threads.length === 0 ? (
        <EmptyState icon={<CheckIcon />} hideContributionLabel>
          <p>You&apos;re all caught up.</p>
        </EmptyState>
      ) : (
        <>
          <ThreadReferenceList
            initialSession={props.initialSession}
            threads={data.threads}
          />
          <PaginationControls
            path="/unread"
            currentPage={data.current_page}
            totalPages={data.total_pages}
            pageSize={data.page_size}
            onClick={handlePageChange}
          />
        </>
      )}
    </LStack>
  );
}
