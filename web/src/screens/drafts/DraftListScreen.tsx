"use client";

import { Unready } from "src/components/site/Unready";

import { NodeCardRows } from "@/components/library/NodeCardList";
import { Heading } from "@/components/ui/heading";
import { VStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../library/useLibraryPath";

import { Props, useDraftListScreen } from "./useDraftListScreen";

export function DraftListScreen(props: Props) {
  const { ready, data, error } = useDraftListScreen(props);
  const libraryPath = useLibraryPath();

  if (!ready) return <Unready error={error} />;

  return (
    <VStack w="full" alignItems="start">
      <Heading>Your drafts</Heading>

      <NodeCardRows
        libraryPath={libraryPath}
        context="generic"
        nodes={data.nodes.data.nodes}
      />
    </VStack>
  );
}
