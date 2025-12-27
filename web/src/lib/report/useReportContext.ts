"use client";

import { useSearchParams } from "next/navigation";

import { handle } from "@/api/client";
import { reportUpdate } from "@/api/openapi-client/reports";
import { ReportStatus } from "@/api/openapi-schema";
import { useSession } from "@/auth";

export function useReportContext() {
  const searchParams = useSearchParams();
  const session = useSession();
  const reportId = searchParams.get("ctx-report-id");

  async function resolveReport() {
    if (!reportId || !session) {
      return;
    }

    await handle(
      async () => {
        await reportUpdate(reportId, {
          status: ReportStatus.resolved,
          handled_by: session.id,
        });
      },
      {
        errorToast: false,
        onError: async (error) => {
          console.error("Failed to resolve report:", error);
        },
      },
    );
  }

  return {
    reportId,
    resolveReport,
  };
}
