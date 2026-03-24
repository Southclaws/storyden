import { z } from "zod";

import { getServerSession } from "@/auth/server-session";
import { accountList } from "@/api/openapi-server/accounts";
import { Permission } from "@/api/openapi-schema";
import { MemberIndexScreen } from "src/screens/library/members/MemberIndexScreen/MemberIndexScreen";
import { AdminMemberIndexScreen } from "src/screens/library/members/AdminMemberIndexScreen/AdminMemberIndexScreen";

import { profileList } from "@/api/openapi-server/profiles";
import { UnreadyBanner } from "@/components/site/Unready";
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
  auth_service: z
    .union([z.string(), z.array(z.string())])
    .transform((v) => (Array.isArray(v) ? v : v.split(",")))
    .optional(),
});
type Query = z.infer<typeof QuerySchema>;

type Props = {
  searchParams: Promise<Query>;
};

export default async function Page(props: Props) {
  try {
    const params = QuerySchema.parse(await props.searchParams);
    const session = await getServerSession();
    const canUseAdminMode = hasPermission(session, Permission.ADMINISTRATOR);
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
        auth_service: params.auth_service,
      });

      return (
        <AdminMemberIndexScreen
          initialResult={data}
          query={params.q}
          page={params.page}
          initialAdmin={params.admin}
          initialSuspended={params.suspended}
          initialAuthServices={params.auth_service}
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
