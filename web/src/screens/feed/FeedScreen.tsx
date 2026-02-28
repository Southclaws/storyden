import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import { categoryListCached } from "@/lib/category/server-category-list";
import { getCategoryThreadListParams } from "@/lib/feed/category";
import { nodeListCached } from "@/lib/library/server-node-list";
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
    <VStack>
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
