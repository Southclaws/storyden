import { linkList } from "@/api/openapi-server/links";
import { nodeList } from "@/api/openapi-server/nodes";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

export default async function Page() {
  const [nodes, links] = await Promise.all([nodeList(), linkList()]);

  return <LibraryIndexScreen nodes={nodes.data} links={links.data} />;
}
