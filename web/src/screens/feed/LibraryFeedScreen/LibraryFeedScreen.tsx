"use client";

import { useNodeGet, useNodeList } from "@/api/openapi-client/nodes";
import {
  type NodeGetOKResponse,
  type NodeListResult,
} from "@/api/openapi-schema";
import { NodeCardGrid, NodeCardRows } from "@/components/library/NodeCardList";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { type FeedBlock } from "@/lib/settings/feed";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";

export type Props = {
  initialNodeList?: NodeListResult;
  initialNode?: NodeGetOKResponse;
  block: Extract<FeedBlock, { type: "library" }>;
};

export function LibraryFeedScreen({
  initialNodeList,
  initialNode,
  block,
}: Props) {
  if (block.node) {
    return <LibraryFeedNode initialNode={initialNode} nodeID={block.node} />;
  }

  return (
    <LibraryFeedRoot
      initialNodeList={initialNodeList}
      layoutType={block.layout}
    />
  );
}

function LibraryFeedNode({
  initialNode,
  nodeID,
}: {
  initialNode?: NodeGetOKResponse;
  nodeID: string;
}) {
  const { data, error } = useNodeGet(
    nodeID,
    {},
    {
      swr: { fallbackData: initialNode },
    },
  );
  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LibraryPageScreen
      embedded
      node={data}
      //childNodes={[]} // TODO: Replicate LibraryPageScreen behavior
    />
  );
}

function LibraryFeedRoot({
  initialNodeList,
  layoutType,
}: {
  initialNodeList?: NodeListResult;
  layoutType: "grid" | "list";
}) {
  const { data, error } = useNodeList(
    {
      //
    },
    {
      swr: { fallbackData: initialNodeList },
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
