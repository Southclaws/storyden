import { Account } from "@/api/openapi-schema";
import { categoryList } from "@/api/openapi-server/categories";
import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";
import { FeedConfig } from "@/components/feed/FeedConfig/FeedConfig";
import { UnreadyBanner } from "@/components/site/Unready";
import { Settings } from "@/lib/settings/settings";
import { VStack } from "@/styled-system/jsx";

import { CategoryIndexScreen } from "../category/CategoryIndexScreen";

import { LibraryFeedScreen } from "./LibraryFeedScreen";
import { ThreadFeedScreen } from "./ThreadFeedScreen/ThreadFeedScreen";

export type PageProps = {
  page: number;
};

export type Props = PageProps & {
  initialSession?: Account;
  initialSettings: Settings;
};

export function FeedScreen({ page, initialSession, initialSettings }: Props) {
  return (
    <VStack>
      <FeedConfig
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
      <FeedScreenContent page={page} initialSettings={initialSettings} />
    </VStack>
  );
}

async function FeedScreenContent({ page, initialSettings }: Props) {
  const feedConfig = initialSettings.metadata.feed;

  switch (feedConfig.source.type) {
    case "threads":
      return <ThreadFeedScreenContent page={page} />;

    case "library":
      return <LibraryFeedScreenContent />;

    case "categories":
      return <CategoryFeedScreenContent />;
  }
}

async function ThreadFeedScreenContent({ page }: PageProps) {
  try {
    const threads = await threadList({
      page: page.toString(),
    });

    return (
      <ThreadFeedScreen initialPage={page} initialPageData={[threads.data]} />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}

async function LibraryFeedScreenContent() {
  try {
    const nodes = await nodeList();
    return <LibraryFeedScreen initialData={nodes.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}

async function CategoryFeedScreenContent() {
  try {
    const categories = await categoryList();
    return <CategoryIndexScreen initialCategoryList={categories.data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
