import { Client } from "src/screens/directory/datagraph/DatagraphIndexScreen";

import { linkList } from "@/api/openapi-server/links";
import { nodeList } from "@/api/openapi-server/nodes";

export default async function Page() {
  const [nodes, links] = await Promise.all([nodeList(), linkList()]);

  return <Client nodes={nodes.data} links={links.data} />;
}
