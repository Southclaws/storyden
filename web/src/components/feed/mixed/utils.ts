import { zip } from "lodash";
import { chunk } from "lodash/fp";

import { NodeList, ThreadList } from "src/api/openapi/schemas";

import { MixedContent, MixedContentLists } from "../useFeed";

export type MixedContentChunk = MixedContentLists & { id: string };

const chunkThreads = chunk(5);
const chunkNodes = chunk(2);

// We chunk the content feed into smaller sections so we're not displaying 100
// threads in one section
export function chunkData(data: MixedContent): MixedContentChunk[] {
  const { threads } = data.threads;
  const nodes = data.nodes?.nodes;

  const threadsChunks = chunkThreads(threads);
  const nodesChunks = chunkNodes(nodes);

  const zipped = zip(threadsChunks, nodesChunks).map(
    (v, i) =>
      ({
        id: i.toString(),
        threads: (v[0] as ThreadList) ?? [],
        nodes: (v[1] as NodeList) ?? [],
      }) satisfies MixedContentChunk,
  );

  return zipped;
}
