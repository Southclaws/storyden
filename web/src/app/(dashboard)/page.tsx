import { UnreadyBanner } from "src/components/site/Unready";

import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";
import { FeedScreen } from "@/screens/feed/FeedScreen";

export default async function Page() {
  try {
    const session = await getServerSession();
    const settings = await getSettings();

    return <FeedScreen initialSession={session} initialSettings={settings} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
