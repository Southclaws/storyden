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
import { ThreadFeedScreen } from "./ThreadFeedScreen";

export type Props = {
  initialSession?: Account;
  initialSettings: Settings;
};

export function FeedScreen({ initialSession, initialSettings }: Props) {
  return (
    <VStack>
      <FeedConfig
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
      <FeedScreenContent initialSettings={initialSettings} />
    </VStack>
  );
}

async function FeedScreenContent({ initialSettings }: Props) {
  const feedConfig = initialSettings.metadata.feed;

  switch (feedConfig.source.type) {
    case "threads":
      return <ThreadFeedScreenContent />;

    case "library":
      return <LibraryFeedScreenContent />;

    case "categories":
      return <CategoryFeedScreenContent />;
  }
}

async function ThreadFeedScreenContent() {
  try {
    const threads = await threadList();
    return <ThreadFeedScreen initialData={threads.data} />;
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
