"use client";

import { Unready } from "src/components/site/Unready";

import { NodeCardRows } from "@/components/library/NodeCardList";
import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { Heading } from "@/components/ui/heading";
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
      <Heading>Your drafts</Heading>

      <Heading color="fg.subtle">Threads</Heading>
      <ThreadReferenceList threads={threads} />

      <Heading color="fg.subtle">Library</Heading>
      <NodeCardRows libraryPath={libraryPath} context="generic" nodes={nodes} />
    </VStack>
  );
}
