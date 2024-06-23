"use client";

import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { NodeCardRows } from "src/components/directory/datagraph/NodeCardList";
import { Unready } from "src/components/site/Unready";

import { DirectoryEmptyState } from "@/components/directory/DirectoryEmptyState";
import { Heading } from "@/components/ui/heading";
import { LStack, VStack } from "@/styled-system/jsx";

import { Props, useDatagraphIndexScreen } from "./useDatagraphIndexScreen";

export function Client(props: Props) {
  const { ready, data, empty, error, session } = useDatagraphIndexScreen(props);

  if (!ready) return <Unready {...error} />;

  const { nodes } = data;

  return (
    <LStack gap="4">
      <Breadcrumbs directoryPath={[]} visibility="draft" create="show" />

      {empty ? (
        <DirectoryEmptyState />
      ) : (
        <p>You can browse the community&apos;s knowledgebase here.</p>
      )}

      {nodes.data.nodes.length > 0 && (
        <NodeCardRows
          directoryPath={[]}
          context="directory"
          size="small"
          {...nodes.data}
        />
      )}
    </LStack>
  );
}
