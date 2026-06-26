import { AdminScreen } from "@/screens/admin/AdminScreen";

import { Permission } from "@/api/openapi-schema";
import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { hasPermission } from "@/utils/permissions";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

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
