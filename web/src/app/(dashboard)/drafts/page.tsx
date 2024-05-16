import {
  NodeListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { DraftListScreen } from "src/screens/drafts/DraftListScreen";

export default async function Page() {
  const [threads, nodes] = await Promise.all([
    server<ThreadListOKResponse>({
      url: "/v1/threads",
      params: { visibility: [Visibility.draft] },
    }),

    server<NodeListOKResponse>({
      url: "/v1/nodes",
      params: { visibility: [Visibility.draft] },
    }),
  ]);

  return <DraftListScreen threads={threads} nodes={nodes} />;
}
