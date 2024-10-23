"use client";

import { Unready } from "src/components/site/Unready";

import { useNodeList } from "@/api/openapi-client/nodes";
import { NodeListResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/feed/EmptyState";
import { NodeCardGrid } from "@/components/library/NodeCardList";

export type Props = {
  initialData?: NodeListResult;
};

export function LibraryFeedScreen({ initialData }: Props) {
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
