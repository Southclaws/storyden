"use client";

import { NodeCardRows } from "@/components/library/NodeCardList";
import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { Unready } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { VStack } from "@/styled-system/jsx";

import { useLibraryPath } from "../library/useLibraryPath";

import { Props, useDraftListScreen } from "./useDraftListScreen";

export function DraftListScreen(props: Props) {
  const { ready, data, error } = useDraftListScreen(props);
  const libraryPath = useLibraryPath();

  if (!ready) return <Unready error={error} />;

  const { nodes, threads } = data;

  return (
    <VStack w="full" alignItems="start">
      <PageHeading>Your drafts</PageHeading>

      <SectionHeading>Threads</SectionHeading>
      <ThreadReferenceList threads={threads} />

      <SectionHeading>Library</SectionHeading>
      <NodeCardRows libraryPath={libraryPath} context="generic" nodes={nodes} />
    </VStack>
  );
}
