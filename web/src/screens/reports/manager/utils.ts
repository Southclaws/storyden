import { DatagraphItemKind, ReportStatus } from "@/api/openapi-schema";

const REPORT_STATUS_VALUES = Object.values(ReportStatus) as ReportStatus[];

export function isReportStatus(
  value: string | null | undefined,
): value is ReportStatus {
  return value ? (REPORT_STATUS_VALUES as string[]).includes(value) : false;
}

export function isDatagraphKind(
  value: string | null | undefined,
): value is DatagraphItemKind {
  return value
    ? (Object.values(DatagraphItemKind) as string[]).includes(value)
    : false;
}
