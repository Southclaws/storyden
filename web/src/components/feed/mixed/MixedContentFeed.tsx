import { PropsWithChildren } from "react";

import { NewBadge } from "src/components/directory/DirectoryBadge";
import { NodeCardGrid } from "src/components/directory/datagraph/NodeCardList";
import { LinkCardRows } from "src/components/directory/links/LinkCardList";

import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, LStack } from "@/styled-system/jsx";

import { TextPostList } from "../text/TextPostList";
import { MixedContent } from "../useFeed";

import { MixedContentChunk, chunkData } from "./utils";

type Props = {
  data: MixedContent;
};

export function MixedContentFeed({ data }: Props) {
  const chunks = chunkData(data);

  return (
    <LStack>
      {chunks.map((c: MixedContentChunk, i) => {
        const showHeader =
          i > 0 ? getChunkType(chunks[i - 1]) === getChunkType(c) : true;

        return (
          <MixedContentFeedSection
            key={c.id}
            data={c}
            showHeader={showHeader}
          />
        );
      })}
    </LStack>
  );
}

// TODO: Really, we should be computing the whole list as a data structure first
// then passing it to a simple mapping component which just performs the render.
function MixedContentFeedSection({
  data,
  showHeader,
}: {
  data: MixedContentChunk;
  showHeader: boolean;
}) {
  const recent = false; // TODO: derive from created/updated dates of the chunks

  return (
    <LStack>
      {data.nodes.length > 0 && (
        <>
          {showHeader && (
            <SectionHeader recent={recent} href="/directory">
              Directory
            </SectionHeader>
          )}
          <NodeCardGrid
            directoryPath={[]}
            context="generic"
            nodes={data.nodes}
          />
        </>
      )}

      {data.threads.length > 0 && (
        <>
          {showHeader && (
            <SectionHeader
              href="/t"
              recent={recent}
              // TODO: Build thread index page then link there.
              // OR even, chunk these based on category and link to the category?
              // href="/t"
            >
              Discussions
            </SectionHeader>
          )}

          {/* TODO: Update this to be a card row list */}
          {/* TODO: Also add the category to the thread IF it's being shown in a context where there are many threads from many categories */}
          <TextPostList posts={data.threads} />
        </>
      )}

      {data.links.length > 0 && (
        <>
          {showHeader && (
            <SectionHeader recent={recent} href="/directory">
              Links
            </SectionHeader>
          )}
          <LinkCardRows links={data.links} />
        </>
      )}
    </LStack>
  );
}

function SectionHeader({
  children,
  recent,
  href,
}: PropsWithChildren<{ recent?: boolean; href?: string }>) {
  return (
    <HStack justify="space-between" w="full">
      <HStack gap="2">
        <Heading>{children}</Heading>
        {recent && <NewBadge />}
      </HStack>

      {href && (
        <LinkButton size="xs" variant="ghost" href={href}>
          See all
        </LinkButton>
      )}
    </HStack>
  );
}

function getChunkType(c: MixedContentChunk | undefined) {
  if (c === undefined) return undefined;
  if (c.nodes.length > 0) return "nodes" as const;
  if (c.threads.length > 0) return "threads" as const;
  if (c.links.length > 0) return "links" as const;
  return undefined;
}
