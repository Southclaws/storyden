import z from "zod";

import {
  DatagraphItemKind,
  Permission,
  ReportStatus,
} from "@/api/openapi-schema";
import { reportList } from "@/api/openapi-server/reports";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { ReportScreen } from "@/screens/reports/ReportScreen";
import { hasPermission } from "@/utils/permissions";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
  status: z.nativeEnum(ReportStatus).optional(),
  kind: z.nativeEnum(DatagraphItemKind).optional(),
});
type Query = z.infer<typeof QuerySchema>;

export default async function Page({ searchParams }: Props) {
  try {
    const { page, status, kind } = QuerySchema.parse(await searchParams);

    const session = await getServerSession();
    const { data } = await reportList({
      page: page?.toString(),
      status,
      kind,
    });

    const canManageReports = hasPermission(
      session,
      Permission.MANAGE_REPORTS,
      Permission.ADMINISTRATOR,
    );

    return (
      <ReportScreen
        initialData={data}
        initialPage={page ?? 1}
        initialStatus={status}
        initialKind={kind}
        canManageReports={canManageReports}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
