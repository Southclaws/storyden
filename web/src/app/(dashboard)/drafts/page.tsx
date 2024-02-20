import {
  ClusterListOKResponse,
  ThreadListOKResponse,
  Visibility,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { DraftListScreen } from "src/screens/drafts/DraftListScreen";

export default async function Page() {
  const [threads, clusters] = await Promise.all([
    server<ThreadListOKResponse>({
      url: "/v1/threads",
      params: { visibility: [Visibility.draft] },
    }),

    server<ClusterListOKResponse>({
      url: "/v1/clusters",
      params: { visibility: [Visibility.draft] },
    }),
  ]);

  console.log(clusters);

  return <DraftListScreen threads={threads} clusters={clusters} />;
}
