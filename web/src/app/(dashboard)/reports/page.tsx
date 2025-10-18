import {
  Account,
  DatagraphItemKind,
  Permission,
  ReportStatus,
} from "@/api/openapi-schema";
import { accountGet } from "@/api/openapi-server/accounts";
import { reportList } from "@/api/openapi-server/reports";
import { UnreadyBanner } from "@/components/site/Unready";
import { ReportScreen } from "@/screens/reports/ReportScreen";

function first(value: string | string[] | undefined): string | undefined {
  if (Array.isArray(value)) {
    return value[0];
  }
  return value;
}

function parsePage(value: string | undefined): number | undefined {
  if (!value) return undefined;
  const parsed = Number.parseInt(value, 10);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : undefined;
}

const REPORT_STATUS_VALUES = Object.values(ReportStatus) as string[];
const DATAGRAPH_KIND_VALUES = Object.values(DatagraphItemKind) as string[];

function isReportStatus(value: string | undefined): value is ReportStatus {
  return value ? REPORT_STATUS_VALUES.includes(value) : false;
}

function isDatagraphKind(value: string | undefined): value is DatagraphItemKind {
  return value ? DATAGRAPH_KIND_VALUES.includes(value) : false;
}

function canManageReports(account: Account | undefined): boolean {
  if (!account) return false;
  if (account.admin) return true;

  const permissions = new Set<string>();
  for (const role of account.roles ?? []) {
    for (const perm of role.permissions ?? []) {
      permissions.add(perm);
    }
  }

  return (
    permissions.has(Permission.ADMINISTRATOR) ||
    permissions.has(Permission.MANAGE_REPORTS)
  );
}

type Props = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function Page({ searchParams }: Props) {
  try {
    const params = await searchParams;

    const pageParam = parsePage(first(params.page));
    const statusParam = first(params.status);
    const kindParam = first(params.kind);

    const status = isReportStatus(statusParam) ? statusParam : undefined;
    const kind = isDatagraphKind(kindParam) ? kindParam : undefined;

    const [reportResponse, accountResponse] = await Promise.all([
      reportList({
        page: pageParam ? pageParam.toString() : undefined,
        status,
        kind,
      }),
      accountGet().catch(() => null),
    ]);

    const account = accountResponse?.data;

    return (
      <ReportScreen
        initialData={reportResponse.data}
        initialPage={pageParam ?? 1}
        initialStatus={status}
        initialKind={kind ?? null}
        canManageReports={canManageReports(account)}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
