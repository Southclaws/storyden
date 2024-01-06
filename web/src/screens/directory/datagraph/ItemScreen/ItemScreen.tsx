"use client";

import { ContentViewer } from "src/components/content/ContentViewer/ContentViewer";
import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { DatagraphHeader } from "src/components/directory/datagraph/Header";
import { Empty } from "src/components/feed/common/PostRef/Empty";
import { Unready } from "src/components/site/Unready";
import { Heading2 } from "src/theme/components/Heading/Index";

import { VStack } from "@/styled-system/jsx";

import { Props, useItemScreen } from "./useItemScreen";

export function ItemScreen(props: Props) {
  const { ready, data, directoryPath, error } = useItemScreen(props);

  if (!ready) return <Unready {...error} />;

  return (
    <VStack w="full" alignItems="start">
      <Breadcrumbs directoryPath={directoryPath} />

      <DatagraphHeader {...props.item} />

      {props.item.content && (
        <VStack>
          <ContentViewer value={props.item.content} />
        </VStack>
      )}

      <VStack alignItems="start" w="full">
        <Heading2>Member of</Heading2>

        {data.clusters.length ? (
          <ClusterList directoryPath={directoryPath} clusters={data.clusters} />
        ) : (
          <Empty>No Items</Empty>
        )}
      </VStack>
    </VStack>
  );
}
