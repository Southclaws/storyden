import { UnreadyBanner } from "@/components/site/Unready";
import { nodeListCached } from "@/lib/library/server-node-list";
import { LibraryIndexScreen } from "@/screens/library/LibraryIndexScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    const nodes = await nodeListCached();

    return <LibraryIndexScreen nodes={nodes.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
