"use client";

import { CreatePageAction } from "@/components/library/CreatePage";
import { LibraryEmptyState } from "@/components/library/LibraryEmptyState";
import { NodeCardGrid } from "@/components/library/NodeCardList";
import { BackAction } from "@/components/site/Action/Back";
import { Unready } from "@/components/site/Unready";
import { PageHeader } from "@/components/ui/page-header";
import { VStack } from "@/styled-system/jsx";

import { Props, useLibraryIndexScreen } from "./useLibraryIndexScreen";

export function LibraryIndexScreen(props: Props) {
  const { ready, data, error } = useLibraryIndexScreen(props);

  if (!ready) return <Unready error={error} />;

  const { nodes } = data;

  return (
    <VStack gap="4">
      <PageHeader
        title="Library"
        back={<BackAction />}
        actions={<CreatePageAction />}
      />

      {nodes.data.nodes.length === 0 ? (
        <LibraryEmptyState />
      ) : (
        <NodeCardGrid libraryPath={[]} context="library" {...nodes.data} />
      )}
    </VStack>
  );
}
