import { nodeGet } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import { categoryListCached } from "@/lib/category/server-category-list";
import { getCategoryThreadListParams } from "@/lib/feed/category";
import { nodeListCached } from "@/lib/library/server-node-list";
import { type FeedBlock } from "@/lib/settings/feed";
import { type FrontendConfiguration } from "@/lib/settings/settings";
import { getSettings } from "@/lib/settings/settings-server";
import { VStack } from "@/styled-system/jsx";

import { FeedScreenContent } from "./FeedScreenContent";
import { InitialData } from "./types";

export type Props = {
  page: number;
};

// NOTE: The FeedScreen (index page) is server hydrated on first load, but the
// admin may configure different sources/layouts which is post-hydration and
// client side. In this case, there is no server side hydrated data available.
// Not a problem just worth pointing out here.
export async function FeedScreen({ page }: Props) {
  const initialSession = await getServerSession();
  const initialSettings = await getSettings();

  const feedConfig = initialSettings.metadata.feed;
  const initialData = await getInitialFeedData(feedConfig, page);

  return (
    <VStack paddingBlockEnd="4">
      <FeedScreenContent
        initialData={initialData}
        initialSettings={initialSettings}
        initialSession={initialSession}
      />
    </VStack>
  );
}

async function getInitialFeedData(
  feedConfig: FrontendConfiguration["feed"],
  page?: number,
): Promise<InitialData> {
  try {
    const categoriesBlock = getBlock(feedConfig.blocks, "categories");
    const threadsBlock = getBlock(feedConfig.blocks, "threads");
    const libraryBlock = getBlock(feedConfig.blocks, "library");

    const [categories, threads, libraryNode, libraryNodes] = await Promise.all([
      categoriesBlock
        ? categoryListCached().then(({ data }) => data)
        : undefined,
      threadsBlock
        ? threadList(
            getCategoryThreadListParams(threadsBlock.source, page),
            feedRequestOptions,
          ).then(({ data }) => data)
        : undefined,
      libraryBlock?.node
        ? nodeGet(libraryBlock.node, undefined, feedRequestOptions).then(
            ({ data }) => data,
          )
        : undefined,
      libraryBlock && !libraryBlock.node
        ? nodeListCached().then(({ data }) => data)
        : undefined,
    ]);

    return {
      initialCategoryList: categories,
      initialPage: page ?? 1,
      initialThreadList: threads,
      initialThreadSource: threadsBlock?.source,
      initialLibraryNode: libraryNode,
      initialLibraryNodeList: libraryNodes,
    };
  } catch (error) {
    // NOTE: Fall back without erroring here, frontend will not be hydrated but
    // it can try the requests again just in case it was a momentary issue. If
    // that fails, we get the standard error handling flow on the frontend. It
    // does mean SSR requests that fail will get a confusing experience but meh.
    return {};
  }
}

const feedRequestOptions = {
  cache: "no-store" as const,
  next: {
    tags: ["feed"],
    revalidate: 0,
  },
};

function getBlock<T extends FeedBlock["type"]>(
  blocks: FeedBlock[],
  type: T,
): Extract<FeedBlock, { type: T }> | undefined {
  return blocks.find(
    (block): block is Extract<FeedBlock, { type: T }> => block.type === type,
  );
}
