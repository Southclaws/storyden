import { UnreadyBanner } from "@/components/site/Unready";
import { QueueScreen } from "@/screens/queue/QueueScreen";

export default async function Page() {
  try {
    return <QueueScreen />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
