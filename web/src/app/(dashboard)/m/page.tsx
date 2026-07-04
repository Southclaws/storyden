import { z } from "zod";

import { Permission } from "@/api/openapi-schema";
import { accountList } from "@/api/openapi-server/accounts";
import { profileList } from "@/api/openapi-server/profiles";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { AdminMemberIndexScreen } from "@/screens/library/members/AdminMemberIndexScreen/AdminMemberIndexScreen";
import { MemberIndexScreen } from "@/screens/library/members/MemberIndexScreen/MemberIndexScreen";
import { hasPermission } from "@/utils/permissions";

const QuerySchema = z.object({
  mode: z.enum(["default", "admin"]).optional(),
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
  roles: z
    .union([z.string(), z.array(z.string())])
    .transform((v) => (Array.isArray(v) ? v : v.split(",")))
    .optional(),
  invited_by: z
    .union([z.string(), z.array(z.string())])
    .transform((v) => (Array.isArray(v) ? v : v.split(",")))
    .optional(),
  joined: z.string().optional(),
  sort: z.string().optional(),
  admin: z
    .enum(["true", "false"])
    .transform((v) => v === "true")
    .optional(),
  suspended: z
    .enum(["true", "false"])
    .transform((v) => v === "true")
    .optional(),
  kind: z.enum(["human", "bot"]).optional(),
});
type Query = z.infer<typeof QuerySchema>;

type Props = {
  searchParams: Promise<Query>;
};

export default async function Page(props: Props) {
  try {
    const params = QuerySchema.parse(await props.searchParams);
    const session = await getServerSession();
    const canUseAdminMode = hasPermission(session, Permission.VIEW_ACCOUNTS);
    const adminMode = canUseAdminMode && params.mode === "admin";

    if (adminMode) {
      const { data } = await accountList({
        q: params.q,
        page: params.page?.toString(),
        roles: params.roles,
        invited_by: params.invited_by,
        joined: params.joined,
        sort: params.sort,
        admin: params.admin,
        suspended: params.suspended,
        kind: params.kind,
      });

      return (
        <AdminMemberIndexScreen
          initialResult={data}
          query={params.q}
          page={params.page}
          initialSort={params.sort}
        />
      );
    }

    const { data } = await profileList({
      q: params.q,
      page: params.page?.toString(),
      roles: params.roles,
      invited_by: params.invited_by,
      joined: params.joined,
      sort: params.sort,
    });

    return (
      <MemberIndexScreen
        initialResult={data}
        query={params.q}
        page={params.page}
        adminModeAvailable={canUseAdminMode}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
