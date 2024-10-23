import { UnreadyBanner } from "src/components/site/Unready";

import { threadList } from "@/api/openapi-server/threads";
import { getSettings } from "@/lib/settings/settings-server";
import { FeedScreen } from "@/screens/feed/FeedScreen";

export default async function Page() {
  try {
    const settings = await getSettings();
    const threads = await threadList();

    return <FeedScreen initialData={threads.data} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
