import { ReportStatus } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import type { ColorPalette } from "@/styled-system/tokens";

const REPORT_STATUS_LABEL: Record<ReportStatus, string> = {
  [ReportStatus.submitted]: "Submitted",
  [ReportStatus.acknowledged]: "Acknowledged",
  [ReportStatus.resolved]: "Resolved",
};

const REPORT_STATUS_COLOR: Record<ReportStatus, ColorPalette> = {
  [ReportStatus.submitted]: "amber",
  [ReportStatus.acknowledged]: "blue",
  [ReportStatus.resolved]: "green",
};

type Props = {
  status: ReportStatus;
};

export function ReportStatusBadge({ status }: Props) {
  return (
    <Badge colorPalette={REPORT_STATUS_COLOR[status]} fontWeight="medium">
      {REPORT_STATUS_LABEL[status]}
    </Badge>
  );
}
