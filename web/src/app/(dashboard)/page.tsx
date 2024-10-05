import { UnreadyBanner } from "src/components/site/Unready";
import { FeedScreenClient } from "src/screens/feed/FeedScreenClient";

import { threadList } from "@/api/openapi-server/threads";

export default async function Page() {
  try {
    const threads = await threadList();

    return <FeedScreenClient initialData={threads.data} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
