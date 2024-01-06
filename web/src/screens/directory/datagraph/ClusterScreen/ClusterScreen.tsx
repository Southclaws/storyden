"use client";

import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { DatagraphHeader } from "src/components/directory/datagraph/Header";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { Empty } from "src/components/feed/common/PostRef/Empty";
import { Unready } from "src/components/site/Unready";

import { VStack } from "@/styled-system/jsx";

import { Props, useClusterScreen } from "./useClusterScreen";

export function ClusterScreen(props: Props) {
  const { ready, data, directoryPath, error } = useClusterScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <VStack w="full" alignItems="start">
      <Breadcrumbs directoryPath={directoryPath} />

      <DatagraphHeader {...props.cluster} />

      {props.cluster.content && (
        <VStack w="full">
          <ContentViewer value={props.cluster.content} />
        </VStack>
      )}

      <VStack alignItems="start" w="full">
        {data.clusters.length === 0 && data.items.length === 0 && (
          <Empty>Nothing inside</Empty>
        )}

        {data.clusters.length > 0 && (
          <ClusterList directoryPath={directoryPath} clusters={data.clusters} />
        )}

        {data.items.length > 0 && (
          <ItemGrid directoryPath={directoryPath} items={data.items} />
        )}
      </VStack>
    </VStack>
  );
}
