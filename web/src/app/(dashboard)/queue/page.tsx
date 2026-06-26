import { UnreadyBanner } from "@/components/site/Unready";
import { QueueScreen } from "@/screens/queue/QueueScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    return <QueueScreen />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
