import { Visibility } from "src/api/openapi-schema";
import { DraftListScreen } from "src/screens/drafts/DraftListScreen";

import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";

export default async function Page() {
  const [threads, nodes] = await Promise.all([
    threadList({
      /* TODO: Visibility param */
    }),

    nodeList({ visibility: [Visibility.draft] }),
  ]);

  return <DraftListScreen threads={threads.data} nodes={nodes.data} />;
}
