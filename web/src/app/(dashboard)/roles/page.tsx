import { roleList } from "@/api/openapi-server/roles";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { RoleScreen } from "@/screens/role/RoleScreen";

export default async function Page() {
  try {
    const session = await getServerSession();
    const { data } = await roleList();

    return <RoleScreen session={session} initialRoles={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
