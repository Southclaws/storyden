"use client";

import { ClusterCardRows } from "src/components/directory/datagraph/ClusterCardList";
import { Unready } from "src/components/site/Unready";
import { Heading1 } from "src/theme/components/Heading/Index";

import { useDirectoryPath } from "../directory/datagraph/useDirectoryPath";

import { VStack } from "@/styled-system/jsx";

import { Props, useDraftListScreen } from "./useDraftListScreen";

export function DraftListScreen(props: Props) {
  const { ready, data, error } = useDraftListScreen(props);
  const directoryPath = useDirectoryPath();

  if (!ready) return <Unready {...error} />;

  return (
    <VStack w="full" alignItems="start">
      <Heading1>Your drafts</Heading1>

      <ClusterCardRows
        directoryPath={directoryPath}
        context="generic"
        clusters={data.clusters.data.clusters}
      />
    </VStack>
  );
}
