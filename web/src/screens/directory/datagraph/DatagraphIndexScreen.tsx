"use client";

import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { LinkCardList } from "src/components/directory/links/LinkCardList";
import { Empty } from "src/components/feed/common/PostRef/Empty";
import { Unready } from "src/components/site/Unready";
import { Heading1, Heading2 } from "src/theme/components/Heading/Index";

import { VStack } from "@/styled-system/jsx";

import { Props, useDatagraphIndexScreen } from "./useDatagraphIndexScreen";

export function Client(props: Props) {
  const { ready, data, empty, error } = useDatagraphIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  const { items, clusters, links } = data;

  if (empty) {
    return <Empty>This community knowledgebase is empty.</Empty>;
  }

  return (
    <VStack w="full" alignItems="start">
      <Heading1>Directory</Heading1>

      <p>You can browse the community&apos;s knowledgebase here.</p>

      {items.data.items.length > 0 && (
        <VStack w="full" alignItems="start">
          <Heading2>New items</Heading2>
          <ItemGrid directoryPath={[]} {...items.data} />
        </VStack>
      )}

      {links.data.results > 0 && (
        <VStack w="full" alignItems="start">
          <Heading2>New links</Heading2>
          <LinkCardList links={links.data} show={3} />
        </VStack>
      )}

      <VStack w="full" alignItems="start">
        <Heading2>Clusters</Heading2>
        <ClusterList directoryPath={[]} {...clusters.data} />
      </VStack>
    </VStack>
  );
}
