import { nodeList } from "@/api/openapi-server/nodes";

import { LibraryNavigationTreeClient } from "./LibraryNavigationTreeClient";

export async function LibraryNavigationTreeServer() {
  try {
    const { data: initialNodeList } = await nodeList({
      // NOTE: This doesn't work due to a bug in Orval.
      // visibility: ["draft", "review", "unlisted", "published"],
    });

    return <LibraryNavigationTreeClient initialNodeList={initialNodeList} />;
  } catch (e) {
    return <LibraryNavigationTreeClient />;
  }
}
