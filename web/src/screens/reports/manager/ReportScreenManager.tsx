"use client";

import { useReportList } from "@/api/openapi-client/reports";
import {
  DatagraphItemKind,
  ReportListResult,
  ReportStatus,
} from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { Center, LStack, styled } from "@/styled-system/jsx";

import { ReportFilters } from "./ReportFilters";
import { ReportResultList } from "./ReportResultList";
import { useReportsScreenFilters } from "./useReportsScreenFilters";

type Props = {
  initialData: ReportListResult;
  initialPage: number;
  initialStatus: ReportStatus | undefined;
  initialKind: DatagraphItemKind | undefined;
  canManageReports: boolean;
};

function useReportsScreen(props: Props) {
  const { page, status, kind } = useReportsScreenFilters();

  const { data, error } = useReportList(
    {
      page: page.toString(),
      status: status ?? undefined,
      kind: kind ?? undefined,
    },
    {
      swr: {
        fallbackData: props.initialData,
        revalidateIfStale: true,
        revalidateOnReconnect: true,
      },
    },
  );

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data,
    page,
    filters: {
      status,
      kind,
    },
  };
}

export function ReportScreenManager(props: Props) {
  const { page, setPage } = useReportsScreenFilters();
  const { ready, error, data, filters } = useReportsScreen(props);

  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const { reports, current_page, total_pages, page_size, results } = data;

  const params: Record<string, string> = {};
  if (filters.status) {
    params["status"] = filters.status;
  }
  if (filters.kind) {
    params["kind"] = filters.kind;
  }

  return (
    <LStack gap="4">
      <Heading>Reports</Heading>

      <ReportFilters />

      {reports.length > 0 ? (
        <>
          <ReportResultList reports={reports} />

          <PaginationControls
            path="/reports"
            params={params}
            currentPage={current_page}
            totalPages={total_pages}
            pageSize={page_size}
            onClick={setPage}
          />
        </>
      ) : (
        <Center w="full">
          <EmptyState hideContributionLabel>No reports to show.</EmptyState>
        </Center>
      )}
    </LStack>
  );
}
