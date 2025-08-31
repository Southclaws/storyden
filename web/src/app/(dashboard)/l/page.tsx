import { nodeList } from "@/api/openapi-server/nodes";
import { UnreadyBanner } from "@/components/site/Unready";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

export default async function Page() {
  try {
    const nodes = await nodeList();

    return <LibraryIndexScreen nodes={nodes.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
