import { ClusterList } from "src/components/directory/datagraph/ClusterList";
import { ItemGrid } from "src/components/directory/datagraph/ItemGrid";
import { LinkListRows } from "src/components/directory/links/LinkListRows";

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
      <ClusterList directoryPath={[]} clusters={data.clusters} />

      <TextPostList posts={data.threads} />

      <ItemGrid directoryPath={[]} items={data.items} />

      <LinkListRows links={data.links} />
    </LStack>
  );
}
