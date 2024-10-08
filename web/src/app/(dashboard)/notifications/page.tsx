import { notificationList } from "@/api/openapi-server/notifications";
import { Unready } from "@/components/site/Unready";
import { NotificationScreen } from "@/screens/notifications/NotificationScreen";

export default async function Page() {
  try {
    const { data } = await notificationList();

    return <NotificationScreen initialData={data} />;
  } catch (e) {
    return <Unready error={e} />;
  }
}
