import { Account } from "@/api/openapi-schema";
import { categoryList } from "@/api/openapi-server/categories";
import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";
import { FeedConfig } from "@/components/feed/FeedConfig/FeedConfig";
import { FrontendConfiguration, Settings } from "@/lib/settings/settings";
import { VStack } from "@/styled-system/jsx";

import { FeedContext } from "./FeedContext";
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
export async function FeedScreen({
  page,
  initialSession,
  initialSettings,
}: Props) {
  const feedConfig = initialSettings.metadata.feed;
  const initialData = await getInitialFeedData(feedConfig, page);

  return (
    <VStack>
      <FeedContext
        initialSession={initialSession}
        initialSettings={initialSettings}
      >
        <FeedConfig />
        <FeedScreenContent initialData={initialData} />
      </FeedContext>
    </VStack>
  );
}

async function getInitialFeedData(
  feedConfig: FrontendConfiguration["feed"],
  page?: number,
): Promise<InitialData> {
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
        library: (await nodeList()).data,
      };

    case "categories":
      return {
        categories: (await categoryList()).data,
      };
  }
}
