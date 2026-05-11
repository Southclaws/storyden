import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { getReportListKey, reportUpdate } from "@/api/openapi-client/reports";
import { Report } from "@/api/openapi-schema";
import {
  DatagraphItemBadge,
  DatagraphItemCard,
} from "@/components/datagraph/DatagraphItemCard";
import { Timestamp } from "@/components/site/Timestamp";
import { Button } from "@/components/ui/button";
import { useI18n } from "@/i18n/provider";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { ReportCardContent } from "../ReportCardContent";
import { ReportStatusBadge } from "../ReportStatusBadge";
import { useReportsScreenFilters } from "../manager/useReportsScreenFilters";

type Props = {
  report: Report;
};

export function ReportCard({ report }: Props) {
  const { mutate } = useSWRConfig();
  const { page } = useReportsScreenFilters();
  const { t } = useI18n();

  async function handleCancel() {
    await handle(
      async () => {
        await reportUpdate(report.id, {
          status: "resolved",
        });
      },
      {
        promiseToast: {
          loading: t("Cancelling report..."),
          success: t("Report cancelled."),
        },
        cleanup: async () => {
          mutate(getReportListKey({ page: page.toString() }));
        },
      },
    );
  }
  return (
    <CardBox w="full" p="2">
      <LStack gap="2" alignItems="stretch">
        <WStack>
          <HStack>
            <ReportStatusBadge status={report.status} />
            <DatagraphItemBadge kind={report.target_kind} />
            <Timestamp created={report.createdAt} large />
          </HStack>

          <HStack>
            <Button
              type="button"
              variant="subtle"
              size="xs"
              disabled={report.status === "resolved"}
              onClick={handleCancel}
            >
              {t("Cancel")}
            </Button>
          </HStack>
        </WStack>

        {report.comment && (
          <styled.blockquote>&ldquo;{report.comment}&rdquo;</styled.blockquote>
        )}

        <ReportCardContent report={report} />
        {/* TODO: Clean up this component and use it in future
        it needs some way to set bg colour maybe? */}
        {/* <DatagraphItemCard item={report.item!} /> */}
      </LStack>
    </CardBox>
  );
}
