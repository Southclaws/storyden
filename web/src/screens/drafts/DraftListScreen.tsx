"use client";

import { NodeCardRows } from "src/components/directory/datagraph/NodeCardList";
import { Unready } from "src/components/site/Unready";

import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";

import { useDirectoryPath } from "../directory/datagraph/useDirectoryPath";

import { Props, useDraftListScreen } from "./useDraftListScreen";

export function DraftListScreen(props: Props) {
  const { ready, data, error } = useDraftListScreen(props);
  const directoryPath = useDirectoryPath();

  if (!ready) return <Unready {...error} />;

  return (
    <VStack w="full" alignItems="start">
      <Heading>Your drafts</Heading>

      <NodeCardRows
        directoryPath={directoryPath}
        context="generic"
        nodes={data.nodes.data.nodes}
      />
    </VStack>
  );
}
