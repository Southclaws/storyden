import { useReportList } from "@/api/openapi-client/reports";
import { ReportListResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { Center, LStack, styled } from "@/styled-system/jsx";

import { useReportsScreenFilters } from "../manager/useReportsScreenFilters";

import { ReportResultList } from "./ReportResultList";

type Props = {
  initialReportsList: ReportListResult;
  initialPage: number;
};

function useReportsScreenMember({ initialReportsList }: Props) {
  const { page } = useReportsScreenFilters();
  const { data, error } = useReportList(
    {
      page: page.toString(),
    },
    {
      swr: {
        fallbackData: initialReportsList,
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
  };
}

export function ReportScreenMember(props: Props) {
  const { ready, error, data, page } = useReportsScreenMember(props);
  if (!ready) {
    return <Unready error={error} />;
  }
  const { reports, results, current_page, total_pages, page_size } = data;
  return (
    <LStack gap="4">
      <Heading>Reports</Heading>

      {reports.length > 0 ? (
        <>
          <ReportResultList reports={reports} />

          <PaginationControls
            path="/reports"
            currentPage={current_page}
            totalPages={total_pages}
            pageSize={page_size}
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
