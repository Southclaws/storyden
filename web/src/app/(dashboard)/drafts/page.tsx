import { Visibility } from "@/api/openapi-schema";
import { DraftListScreen } from "@/screens/drafts/DraftListScreen";

import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import {
  UnauthenticatedBanner,
  UnreadyBanner,
} from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export default async function Page() {
  try {
    const [session, settings] = await Promise.all([
      getServerSession(),
      getSettings(),
    ]);

    if (!session) {
      return <UnauthenticatedBanner initialSettings={settings} />;
    }

    const [threads, nodes] = await Promise.all([
      threadList(
        { author: session.handle, visibility: [Visibility.draft] },
        {
          cache: "no-store",
          next: { revalidate: 0 },
        },
      ),
      nodeList(
        { author: session.handle, visibility: [Visibility.draft] },
        {
          cache: "no-store",
          next: { revalidate: 0 },
        },
      ),
    ]);

    return (
      <DraftListScreen
        session={session}
        initialThreads={threads.data}
        initialNodes={nodes.data}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
