import { zip } from "lodash";
import { chunk, filter } from "lodash/fp";

import { Link, LinkList, NodeList, ThreadList } from "src/api/openapi/schemas";

import { MixedContent, MixedContentLists } from "../useFeed";

export type MixedContentChunk = MixedContentLists & { id: string };

const chunkThreads = chunk(5);
const chunkNodes = chunk(2);
const chunkLinks = chunk(5);

const filterInterestingLinks = filter((v: Link) =>
  Boolean(v.title && v.description && v.assets.length > 0),
);

// We chunk the content feed into smaller sections so we're not displaying 100
// threads in one section
export function chunkData(data: MixedContent): MixedContentChunk[] {
  const { threads } = data.threads;
  const { nodes } = data.nodes;
  const { links } = data.links;

  // Filter out links that are missing titles/etc so the home screen looks nice.
  const goodLinks = filterInterestingLinks(links);

  const threadsChunks = chunkThreads(threads);
  const nodesChunks = chunkNodes(nodes);
  const linksChunks = chunkLinks(goodLinks);

  const zipped = zip(threadsChunks, nodesChunks, linksChunks).map(
    (v, i) =>
      ({
        id: i.toString(),
        threads: (v[0] as ThreadList) ?? [],
        nodes: (v[1] as NodeList) ?? [],
        links: (v[3] as LinkList) ?? [],
      }) satisfies MixedContentChunk,
  );

  return zipped;
}
