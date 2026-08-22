import { ReportStatus } from "@/api/openapi-schema";
import {
  StatusBadge,
  type StatusBadgeTone,
} from "@/components/ui/status-badge";

const REPORT_STATUS_LABEL: Record<ReportStatus, string> = {
  [ReportStatus.submitted]: "Submitted",
  [ReportStatus.acknowledged]: "Acknowledged",
  [ReportStatus.resolved]: "Resolved",
};

const REPORT_STATUS_TONE: Record<ReportStatus, StatusBadgeTone> = {
  [ReportStatus.submitted]: "warning",
  [ReportStatus.acknowledged]: "info",
  [ReportStatus.resolved]: "success",
};

type Props = {
  status: ReportStatus;
};

export function ReportStatusBadge({ status }: Props) {
  return (
    <StatusBadge tone={REPORT_STATUS_TONE[status]} fontWeight="medium">
      {REPORT_STATUS_LABEL[status]}
    </StatusBadge>
  );
}
