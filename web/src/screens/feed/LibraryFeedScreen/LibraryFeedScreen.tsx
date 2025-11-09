"use client";

import { Unready } from "src/components/site/Unready";

import { useNodeGet, useNodeList } from "@/api/openapi-client/nodes";
import { NodeListResult } from "@/api/openapi-schema";
import { NodeCardGrid, NodeCardRows } from "@/components/library/NodeCardList";
import { EmptyState } from "@/components/site/EmptyState";
import { FeedConfig } from "@/lib/settings/feed";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";

export type Props = {
  initialData?: NodeListResult;
  initialSettings: Settings;
};

export function LibraryFeedScreen({ initialData, initialSettings }: Props) {
  const { settings } = useSettings(initialSettings);
  const feed = settings?.metadata.feed ?? initialSettings.metadata.feed;

  if (feed.source.type !== "library") {
    return null;
  }

  if (feed.source.node) {
    return (
      <LibraryFeedNode
        initialSettings={initialSettings}
        initialData={initialData}
        nodeID={feed.source.node}
      />
    );
  }

  return (
    <LibraryFeedRoot
      initialSettings={initialSettings}
      initialData={initialData}
      feed={feed}
    />
  );
}

function LibraryFeedNode({ initialData, nodeID }: Props & { nodeID: string }) {
  const { data, error } = useNodeGet(
    nodeID,
    {},
    {
      // swr: { fallbackData: initialData },
    },
  );
  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LibraryPageScreen
      node={data}
      //childNodes={[]} // TODO: Replicate LibraryPageScreen behavior
    />
  );
}

function LibraryFeedRoot({ initialData, feed }: Props & { feed: FeedConfig }) {
  const { data, error } = useNodeList(
    {
      //
    },
    {
      swr: { fallbackData: initialData },
    },
  );
  if (!data) {
    return <Unready error={error} />;
  }

  if (data.nodes.length === 0) {
    return <EmptyState />;
  }

  switch (feed.layout.type) {
    case "grid":
      return (
        <NodeCardGrid libraryPath={[]} nodes={data.nodes} context="library" />
      );

    case "list":
      return (
        <NodeCardRows libraryPath={[]} nodes={data.nodes} context="library" />
      );
  }
}
