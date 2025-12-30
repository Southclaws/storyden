import { UnreadyBanner } from "@/components/site/Unready";
import { nodeListCached } from "@/lib/library/server-node-list";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

export default async function Page() {
  try {
    const nodes = await nodeListCached();

    return <LibraryIndexScreen nodes={nodes.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
