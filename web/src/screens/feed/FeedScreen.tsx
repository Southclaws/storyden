import { Account } from "@/api/openapi-schema";
import { threadList } from "@/api/openapi-server/threads";
import { categoryListCached } from "@/lib/category/server-category-list";
import { getCategoryThreadListParams } from "@/lib/feed/category";
import { nodeListCached } from "@/lib/library/server-node-list";
import { FrontendConfiguration, Settings } from "@/lib/settings/settings";
import { VStack } from "@/styled-system/jsx";

import { FeedScreenContent, InitialData } from "./FeedScreenContent";

export type PageProps = {
  initialSession?: Account;
  page: number;
};

export type Props = PageProps & {
  initialSettings: Settings;
};

// NOTE: The FeedScreen (index page) is server hydrated on first load, but the
// admin may configure different sources/layouts which is post-hydration and
// client side. In this case, there is no server side hydrated data available.
// Not a problem just worth pointing out here.
export async function FeedScreen({ page, initialSettings }: Props) {
  const feedConfig = initialSettings.metadata.feed;
  const initialData = await getInitialFeedData(feedConfig, page);

  return (
    <VStack>
      <FeedScreenContent initialData={initialData} />
    </VStack>
  );
}

async function getInitialFeedData(
  feedConfig: FrontendConfiguration["feed"],
  page?: number,
): Promise<InitialData> {
  try {
    switch (feedConfig.source.type) {
      case "threads":
        return {
          page: page ?? 1,
          threads: (
            await threadList(
              {
                page: page?.toString(),
              },
              {
                cache: "no-store",
                next: {
                  tags: ["feed"],
                  revalidate: 0,
                },
              },
            )
          ).data,
        };

      case "library":
        return {
          library: (await nodeListCached()).data,
        };

      case "categories": {
        const mode = feedConfig.source.threadListMode ?? "all";
        const categories = (await categoryListCached()).data;

        const threadParams = getCategoryThreadListParams(mode, page);
        const threads = (
          await threadList(threadParams, {
            cache: "no-store",
            next: {
              tags: ["feed"],
              revalidate: 0,
            },
          })
        ).data;

        return {
          categories,
          page: page ?? 1,
          threads,
        };
      }
    }
  } catch (error) {
    // NOTE: Fall back without erroring here, frontend will not be hydrated but
    // it can try the requests again just in case it was a momentary issue. If
    // that fails, we get the standard error handling flow on the frontend. It
    // does mean SSR requests that fail will get a confusing experience but meh.
    return {};
  }
}
