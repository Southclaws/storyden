"use client";

import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { NodeCardRows } from "src/components/directory/datagraph/NodeCardList";
import { LinkResultList } from "src/components/directory/links/LinkIndexView/LinkResultList";
import { Empty } from "src/components/site/Empty";
import { Unready } from "src/components/site/Unready";

import { Heading2 } from "@/components/ui/typography-heading";
import { Center, VStack } from "@/styled-system/jsx";

import { Props, useDatagraphIndexScreen } from "./useDatagraphIndexScreen";

export function Client(props: Props) {
  const { ready, data, empty, error, session } = useDatagraphIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  const { nodes, links } = data;

  return (
    <VStack w="full" alignItems="start" gap="4">
      <Breadcrumbs directoryPath={[]} visibility="draft" create="show" />

      {empty ? (
        <Center h="full">
          <Empty>
            This community knowledgebase is empty.
            <br />
            {session ? (
              <>Be the first to contribute!</>
            ) : (
              <>Please log in to contribute.</>
            )}
          </Empty>
        </Center>
      ) : (
        <p>You can browse the community&apos;s knowledgebase here.</p>
      )}

      {links.data.results > 0 && (
        <VStack w="full" alignItems="start">
          <Heading2>New links</Heading2>
          <LinkResultList links={links.data} show={3} />
        </VStack>
      )}

      {nodes.data.nodes.length > 0 && (
        <VStack w="full" alignItems="start">
          <Heading2>Knowledgebase</Heading2>
          <NodeCardRows
            directoryPath={[]}
            context="directory"
            size="small"
            {...nodes.data}
          />
        </VStack>
      )}
    </VStack>
  );
}
