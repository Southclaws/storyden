"use client";

import {
  DatagraphItemKind,
  ReportListResult,
  ReportStatus,
} from "@/api/openapi-schema";

import { ReportScreenManager } from "./manager/ReportScreenManager";
import { ReportScreenMember } from "./member/ReportScreenMember";

type Props = {
  initialData: ReportListResult;
  initialPage: number;
  initialStatus: ReportStatus | undefined;
  initialKind: DatagraphItemKind | undefined;
  canManageReports: boolean;
};

export function ReportScreen(props: Props) {
  switch (props.canManageReports) {
    case true:
      return <ReportScreenManager {...props} />;
    case false:
      return (
        <ReportScreenMember
          initialReportsList={props.initialData}
          initialPage={props.initialPage}
        />
      );
  }
}
