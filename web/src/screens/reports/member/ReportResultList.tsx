import { ReportList } from "@/api/openapi-schema";
import { LStack } from "@/styled-system/jsx";

import { ReportCard } from "./ReportCard";

type Props = {
  reports: ReportList;
};

export function ReportResultList({ reports }: Props) {
  return (
    <LStack gap="3">
      {reports.map((report) => (
        <ReportCard key={report.id} report={report} />
      ))}
    </LStack>
  );
}
