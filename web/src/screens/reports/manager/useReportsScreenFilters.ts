import { parseAsInteger, useQueryState } from "nuqs";

import { DatagraphItemKind, ReportStatus } from "@/api/openapi-schema";

import { isDatagraphKind, isReportStatus } from "./utils";

export function useReportsScreenFilters() {
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(1));

  const [status, setStatus] = useQueryState<ReportStatus | null>("status", {
    defaultValue: null,
    clearOnDefault: true,
    parse(value) {
      if (!value) return null;
      return isReportStatus(value) ? (value as ReportStatus) : null;
    },
  });

  const [kind, setKind] = useQueryState<DatagraphItemKind | null>("kind", {
    defaultValue: null,
    clearOnDefault: true,
    parse(value) {
      if (!value) return null;
      return isDatagraphKind(value) ? (value as DatagraphItemKind) : null;
    },
  });

  return {
    page,
    setPage,
    status,
    setStatus,
    kind,
    setKind,
  };
}
