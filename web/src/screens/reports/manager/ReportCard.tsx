import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { getReportListKey, reportUpdate } from "@/api/openapi-client/reports";
import { Report, ReportStatus } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { DatagraphItemBadge } from "@/components/datagraph/DatagraphItemCard";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { SystemBadge } from "@/components/system/SystemBadge";
import { Button } from "@/components/ui/button";
import {
  Box,
  CardBox,
  HStack,
  LStack,
  WStack,
  styled,
} from "@/styled-system/jsx";

import { ReportCardContent } from "../ReportCardContent";
import { ReportStatusBadge } from "../ReportStatusBadge";

import { useReportsScreenFilters } from "./useReportsScreenFilters";

type Props = {
  report: Report;
};

export function ReportCard({ report }: Props) {
  const session = useSession();
  const { mutate } = useSWRConfig();
  const { page, status, kind } = useReportsScreenFilters();

  async function handleAcknowledge() {
    await handle(
      async () => {
        await reportUpdate(report.id, {
          status: ReportStatus.acknowledged,
          handled_by: session?.id,
        });
      },
      {
        promiseToast: {
          loading: "Acknowledging report...",
          success: "Report acknowledged.",
        },
        cleanup: async () => {
          mutate(
            getReportListKey({
              page: page.toString(),
              status: status ?? undefined,
              kind: kind ?? undefined,
            }),
          );
        },
      },
    );
  }

  async function handleResolve() {
    await handle(
      async () => {
        await reportUpdate(report.id, {
          status: ReportStatus.resolved,
          handled_by: session?.id,
        });
      },
      {
        promiseToast: {
          loading: "Resolving report...",
          success: "Report resolved.",
        },
        cleanup: async () => {
          mutate(
            getReportListKey({
              page: page.toString(),
              status: status ?? undefined,
              kind: kind ?? undefined,
            }),
          );
        },
      },
    );
  }

  const isAcknowledged = report.status === ReportStatus.acknowledged;
  const isResolved = report.status === ReportStatus.resolved;

  return (
    <CardBox w="full" p="2">
      <LStack gap="3" alignItems="stretch">
        <WStack>
          <HStack>
            <ReportStatusBadge status={report.status} />
            <DatagraphItemBadge kind={report.target_kind} />
          </HStack>
          <Timestamp created={report.createdAt} large />
        </WStack>

        {report.comment && (
          <styled.blockquote>&ldquo;{report.comment}&rdquo;</styled.blockquote>
        )}

        <WStack gap="1">
          <HStack gap="2" alignItems="center" minW="0" maxW="1/2">
            <styled.span fontSize="sm" color="fg.subtle" fontWeight="medium">
              Reporter:
            </styled.span>
            {report.reported_by ? (
              <MemberBadge
                profile={report.reported_by}
                name="handle"
                size="sm"
                as="link"
              />
            ) : (
              <SystemBadge size="sm" name="visible" />
            )}
          </HStack>

          <HStack gap="2" alignItems="center" minW="0" maxW="1/2">
            <styled.span fontSize="sm" color="fg.subtle" fontWeight="medium">
              Handler:
            </styled.span>
            {report.handled_by ? (
              <MemberBadge
                profile={report.handled_by}
                name="handle"
                size="sm"
                as="link"
              />
            ) : (
              <styled.span color="fg.subtle" fontStyle="italic">
                Unassigned
              </styled.span>
            )}
          </HStack>
        </WStack>

        <ReportCardContent report={report} />

        <WStack>
          <Button
            type="button"
            size="xs"
            variant="subtle"
            disabled={isAcknowledged || isResolved}
            onClick={handleAcknowledge}
          >
            Acknowledge
          </Button>
          <Button
            type="button"
            size="xs"
            variant="solid"
            disabled={isResolved}
            onClick={handleResolve}
          >
            Resolve
          </Button>
        </WStack>
      </LStack>
    </CardBox>
  );
}
