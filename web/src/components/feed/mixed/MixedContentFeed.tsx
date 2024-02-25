import { ClusterCardRows } from "src/components/directory/datagraph/ClusterCardList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { LinkCardRows } from "src/components/directory/links/LinkCardList";

import { TextPostList } from "../text/TextPostList";
import { MixedContent } from "../useFeed";

import { LStack } from "@/styled-system/jsx";

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

function MixedContentFeedSection({ data }: { data: MixedContentChunk }) {
  return (
    <LStack>
      <ClusterCardRows directoryPath={[]} clusters={data.clusters} />

      <TextPostList posts={data.threads} />

      <ItemGrid directoryPath={[]} items={data.items} />

      <LinkCardRows links={data.links} />
    </LStack>
  );
}
