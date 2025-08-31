"use client";

import { Unready } from "src/components/site/Unready";

import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryEmptyState } from "@/components/library/LibraryEmptyState";
import { NodeCardGrid } from "@/components/library/NodeCardList";
import { VStack } from "@/styled-system/jsx";

import { Props, useLibraryIndexScreen } from "./useLibraryIndexScreen";

export function LibraryIndexScreen(props: Props) {
  const { ready, data, error } = useLibraryIndexScreen(props);

  if (!ready) return <Unready error={error} />;

  const { nodes } = data;

  return (
    <VStack gap="4">
      <Breadcrumbs libraryPath={[]} visibility="draft" create="show" />

      {nodes.data.nodes.length === 0 ? (
        <LibraryEmptyState />
      ) : (
        <NodeCardGrid libraryPath={[]} context="library" {...nodes.data} />
      )}
    </VStack>
  );
}
