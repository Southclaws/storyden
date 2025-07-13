"use client";

import { Unready } from "src/components/site/Unready";

import { useNodeGet, useNodeList } from "@/api/openapi-client/nodes";
import { NodeListResult } from "@/api/openapi-schema";
import { NodeCardGrid } from "@/components/library/NodeCardList";
import { EmptyState } from "@/components/site/EmptyState";
import { LibraryPageScreen } from "@/screens/library/LibraryPageScreen/LibraryPageScreen";

import { useFeedContext } from "../FeedContext";

export type Props = {
  initialData?: NodeListResult;
};

export function LibraryFeedScreen({ initialData }: Props) {
  const { feed } = useFeedContext();
  if (feed.source.type !== "library") {
    return null;
  }

  if (feed.source.node) {
    return (
      <LibraryFeedNode initialData={initialData} nodeID={feed.source.node} />
    );
  }

  return <LibraryFeedRoot initialData={initialData} />;
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

function LibraryFeedRoot({ initialData }: Props) {
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

  return <NodeCardGrid libraryPath={[]} nodes={data.nodes} context="library" />;
}
