"use client";

import { Unready } from "src/components/site/Unready";

import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { LibraryEmptyState } from "@/components/library/LibraryEmptyState";
import { NodeCardRows } from "@/components/library/NodeCardList";
import { LStack } from "@/styled-system/jsx";

import { Props, useLibraryIndexScreen } from "./useLibraryIndexScreen";

export function LibraryIndexScreen(props: Props) {
  const { ready, data, empty, error } = useLibraryIndexScreen(props);

  if (!ready) return <Unready error={error} />;

  const { nodes } = data;

  return (
    <LStack gap="4">
      <Breadcrumbs libraryPath={[]} visibility="draft" create="show" />

      {empty ? (
        <LibraryEmptyState />
      ) : (
        <p>You can browse the community&apos;s library here.</p>
      )}

      {nodes.data.nodes.length > 0 && (
        <NodeCardRows libraryPath={[]} context="library" {...nodes.data} />
      )}
    </LStack>
  );
}
