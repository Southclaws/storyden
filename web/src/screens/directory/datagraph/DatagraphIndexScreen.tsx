"use client";

import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { Unready } from "src/components/site/Unready";
import { Heading1, Heading2 } from "src/theme/components/Heading/Index";

import { VStack } from "@/styled-system/jsx";

import { Props, useDatagraphIndexScreen } from "./useDatagraphIndexScreen";

export function Client(props: Props) {
  const { ready, data, error } = useDatagraphIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  const { items, clusters } = data;

  return (
    <VStack w="full" alignItems="start">
      <Heading1>Directory</Heading1>

      <p>
        The place where stuff becomes meaningful. Where order and chaos are
        balanced. Where knowledge and data are one and the same. Go forth and
        accumulate, learn and prosper.
      </p>

      <VStack w="full" alignItems="start">
        <Heading2>New</Heading2>
        <ItemGrid directoryPath={[]} {...items.data} />
      </VStack>

      <VStack w="full" alignItems="start">
        <Heading2>Clusters</Heading2>
        <ClusterList directoryPath={[]} {...clusters.data} />
      </VStack>
    </VStack>
  );
}
