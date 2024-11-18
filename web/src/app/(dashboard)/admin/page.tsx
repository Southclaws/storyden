import { AdminScreen } from "src/screens/admin/AdminScreen";

import { Permission } from "@/api/openapi-schema";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { hasPermission } from "@/utils/permissions";

export default async function Page() {
  try {
    const session = await getServerSession();
    if (
      !session ||
      !hasPermission(
        session,
        Permission.ADMINISTRATOR,
        Permission.MANAGE_SETTINGS,
      )
    ) {
      return (
        <UnreadyBanner error="Not authorised to view the system configuration page." />
      );
    }

    return <AdminScreen />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
