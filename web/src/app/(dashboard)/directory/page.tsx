import {
  LinkListOKResponse,
  NodeListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { Client } from "src/screens/directory/datagraph/DatagraphIndexScreen";

export default async function Page() {
  const [nodes, links] = await Promise.all([
    server<NodeListOKResponse>({ url: "/v1/nodes" }),
    server<LinkListOKResponse>({ url: "/v1/links" }),
  ]);

  return <Client nodes={nodes} links={links} />;
}
