import { PropsWithChildren } from "react";

import { NewBadge } from "src/components/directory/DirectoryBadge";
import { ClusterCardGrid } from "src/components/directory/datagraph/ClusterCardList";
import { ItemCardGrid } from "src/components/directory/datagraph/ItemCardList";
import { LinkCardRows } from "src/components/directory/links/LinkCardList";
import { Heading2 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { TextPostList } from "../text/TextPostList";
import { MixedContent } from "../useFeed";

import { HStack, LStack } from "@/styled-system/jsx";

import { MixedContentChunk, chunkData } from "./utils";

type Props = {
  data: MixedContent;
};

export function MixedContentFeed({ data }: Props) {
  const chunks = chunkData(data);

  return (
    <LStack>
      {chunks.map((c: MixedContentChunk) => (
        <MixedContentFeedSection key={c.id} data={c} />
      ))}
    </LStack>
  );
}

// TODO: Really, we should be computing the whole list as a data structure first
// then passing it to a simple mapping component which just performs the render.
function MixedContentFeedSection({ data }: { data: MixedContentChunk }) {
  const recent = true; // TODO: derive from created/updated dates of the chunks

  const dontShowDirectoryHeaderTwice = Boolean(
    !(data.clusters.length && data.items.length && !data.threads.length),
  );

  return (
    <LStack>
      {data.clusters.length > 0 && (
        <>
          <SectionHeader recent={recent} href="/directory">
            Directory
          </SectionHeader>
          <ClusterCardGrid
            directoryPath={[]}
            context="generic"
            clusters={data.clusters}
          />
        </>
      )}

      {data.threads.length > 0 && (
        <>
          <SectionHeader
            href="/t"
            recent={recent}
            // TODO: Build thread index page then link there.
            // OR even, chunk these based on category and link to the category?
            // href="/t"
          >
            Discussions
          </SectionHeader>

          {/* TODO: Update this to be a card row list */}
          {/* TODO: Also add the category to the thread IF it's being shown in a context where there are many threads from many categories */}
          <TextPostList posts={data.threads} />
        </>
      )}

      {data.items.length > 0 && (
        <>
          {dontShowDirectoryHeaderTwice && (
            <SectionHeader recent={recent} href="/directory">
              Directory
            </SectionHeader>
          )}
          <ItemCardGrid directoryPath={[]} items={data.items} />
        </>
      )}

      {data.links.length > 0 && (
        <>
          <SectionHeader recent={recent} href="/directory">
            Links
          </SectionHeader>
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
        <Heading2 size="xs">{children}</Heading2>
        {recent && <NewBadge />}
      </HStack>

      {href && (
        <Link size="xs" href={href}>
          See all
        </Link>
      )}
    </HStack>
  );
}
