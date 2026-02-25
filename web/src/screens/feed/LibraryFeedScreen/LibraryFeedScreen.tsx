"use client";

import { Unready } from "src/components/site/Unready";

import { useNodeGet, useNodeList } from "@/api/openapi-client/nodes";
import { type NodeListResult } from "@/api/openapi-schema";
import { NodeCardGrid, NodeCardRows } from "@/components/library/NodeCardList";
import { EmptyState } from "@/components/site/EmptyState";
import { type FeedConfig } from "@/lib/settings/feed";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";

export type Props = {
  initialData?: NodeListResult;
  feed: FeedConfig;
};

export function LibraryFeedScreen({ initialData, feed }: Props) {
  if (feed.source.type !== "library") {
    return null;
  }

  if (feed.source.node) {
    return (
      <LibraryFeedNode initialData={initialData} nodeID={feed.source.node} />
    );
  }

  return (
    <LibraryFeedRoot initialData={initialData} layoutType={feed.layout.type} />
  );
}

function LibraryFeedNode({
  initialData,
  nodeID,
}: {
  initialData?: NodeListResult;
  nodeID: string;
}) {
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

function LibraryFeedRoot({
  initialData,
  layoutType,
}: {
  initialData?: NodeListResult;
  layoutType: FeedConfig["layout"]["type"];
}) {
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

  switch (layoutType) {
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
