import {
  ClusterListOKResponse,
  LinkListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { Client } from "src/screens/directory/datagraph/DatagraphIndexScreen";

export default async function Page() {
  const [clusters, links] = await Promise.all([
    server<ClusterListOKResponse>({ url: "/v1/clusters" }),
    server<LinkListOKResponse>({ url: "/v1/links" }),
  ]);

  return <Client clusters={clusters} links={links} />;
}
