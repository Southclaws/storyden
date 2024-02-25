import { zip } from "lodash";
import { chunk, filter } from "lodash/fp";

import {
  ClusterList,
  ItemList,
  Link,
  LinkList,
  ThreadList,
} from "src/api/openapi/schemas";

import { MixedContent, MixedContentLists } from "../useFeed";

export type MixedContentChunk = MixedContentLists & { id: string };

const chunkThreads = chunk(5);
const chunkClusters = chunk(2);
const chunkItems = chunk(2);
const chunkLinks = chunk(5);

const filterInterestingLinks = filter((v: Link) =>
  Boolean(v.title && v.description && v.assets.length > 0),
);

// We chunk the content feed into smaller sections so we're not displaying 100
// threads in one section
export function chunkData(data: MixedContent): MixedContentChunk[] {
  const { threads } = data.threads;
  const { clusters } = data.clusters;
  const { items } = data.items;
  const { links } = data.links;

  // Filter out links that are missing titles/etc so the home screen looks nice.
  const goodLinks = filterInterestingLinks(links);

  const threadsChunks = chunkThreads(threads);
  const clustersChunks = chunkClusters(clusters);
  const itemsChunks = chunkItems(items);
  const linksChunks = chunkLinks(goodLinks);

  const zipped = zip(
    threadsChunks,
    clustersChunks,
    itemsChunks,
    linksChunks,
  ).map(
    (v, i) =>
      ({
        id: i.toString(),
        threads: (v[0] as ThreadList) ?? [],
        clusters: (v[1] as ClusterList) ?? [],
        items: (v[2] as ItemList) ?? [],
        links: (v[3] as LinkList) ?? [],
      }) satisfies MixedContentChunk,
  );

  return zipped;
}
