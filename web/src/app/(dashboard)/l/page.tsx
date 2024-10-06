import { linkList } from "@/api/openapi-server/links";
import { nodeList } from "@/api/openapi-server/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

export default async function Page() {
  try {
    const [nodes, links] = await Promise.all([nodeList(), linkList()]);

    return <LibraryIndexScreen nodes={nodes.data} links={links.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
