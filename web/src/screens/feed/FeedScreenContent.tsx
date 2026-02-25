"use client";

import { type Account } from "@/api/openapi-schema";
import { useFeedConfig } from "@/lib/settings/feed-client";
import { type Settings } from "@/lib/settings/settings";

import { CategoryIndexScreen } from "../category/CategoryIndexScreen";

import { LibraryFeedScreen } from "./LibraryFeedScreen/LibraryFeedScreen";
import { ThreadFeedScreen } from "./ThreadFeedScreen/ThreadFeedScreen";
import { InitialData } from "./types";

type Props = {
  initialData: InitialData;
  initialSettings?: Settings;
  initialSession?: Account;
};

export function FeedScreenContent({
  initialData,
  initialSettings,
  initialSession,
}: Props) {
  const feed = useFeedConfig(initialSettings);

  switch (feed.source.type) {
    case "threads":
      return (
        <ThreadFeedScreen
          initialPage={initialData.page}
          initialPageData={initialData.threads}
          initialSession={initialSession}
          initialSettings={initialSettings}
          category={undefined}
          paginationBasePath="/"
          showCategorySelect={true}
          showQuickShare={feed.source.quickShare === "enabled"}
        />
      );

    case "library":
      return (
        <LibraryFeedScreen
          initialData={initialData.library}
          initialSettings={initialSettings}
        />
      );

    case "categories":
      return (
        <CategoryIndexScreen
          layout={feed.layout.type}
          threadListMode={feed.source.threadListMode}
          showQuickShare={feed.source.quickShare === "enabled"}
          initialCategoryList={initialData.categories}
          initialThreadList={initialData.threads}
          initialThreadListPage={initialData.page}
          paginationBasePath="/"
        />
      );
  }
}
