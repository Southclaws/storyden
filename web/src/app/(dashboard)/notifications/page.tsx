import { notificationList } from "@/api/openapi-server/notifications";
import { Unready } from "@/components/site/Unready";
import { NotificationScreen } from "@/screens/notifications/NotificationScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    const { data } = await notificationList();

    return <NotificationScreen initialData={data} />;
  } catch (e) {
    return <Unready error={e} />;
  }
}
