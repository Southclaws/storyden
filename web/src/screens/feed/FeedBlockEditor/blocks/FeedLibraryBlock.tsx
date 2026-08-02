import { useNodeGet, useNodeList } from "@/api/openapi-client/nodes";
import { NodeCardGrid, NodeCardRows } from "@/components/library/NodeCardList";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { FeedBlock } from "@/lib/settings/feed";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";

import { useFeedBlockEditor } from "../Context";

export function FeedLibraryBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "library" }>;
}) {
  if (block.node) {
    return <FeedLibraryPage nodeID={block.node} />;
  }

  return <FeedLibraryRoot layout={block.layout} />;
}

function FeedLibraryPage({ nodeID }: { nodeID: string }) {
  const { initialData } = useFeedBlockEditor();
  const { data, error } = useNodeGet(nodeID, undefined, {
    swr: { fallbackData: initialData.initialLibraryNode },
  });

  if (!data) {
    return <Unready error={error} />;
  }

  return <LibraryPageScreen embedded node={data} />;
}

function FeedLibraryRoot({ layout }: { layout: "grid" | "list" }) {
  const { initialData } = useFeedBlockEditor();
  const { data, error } = useNodeList(undefined, {
    swr: { fallbackData: initialData.initialLibraryNodeList },
  });

  if (!data) {
    return <Unready error={error} />;
  }

  if (data.nodes.length === 0) {
    return <EmptyState />;
  }

  return layout === "grid" ? (
    <NodeCardGrid libraryPath={[]} nodes={data.nodes} context="library" />
  ) : (
    <NodeCardRows libraryPath={[]} nodes={data.nodes} context="library" />
  );
}
