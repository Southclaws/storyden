"use client";

import {
  Account,
  CategoryListResult,
  NodeListResult,
  ThreadListResult,
} from "@/api/openapi-schema";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";

import { CategoryIndexScreen } from "../category/CategoryIndexScreen";

import { LibraryFeedScreen } from "./LibraryFeedScreen/LibraryFeedScreen";
import { ThreadFeedScreen } from "./ThreadFeedScreen/ThreadFeedScreen";

export type InitialData = {
  threads?: ThreadListResult;
  page?: number;
  library?: NodeListResult;
  categories?: CategoryListResult;
};

type Props = {
  initialData: InitialData;
  initialSettings: Settings;
  initialSession?: Account;
};

export function FeedScreenContent({
  initialData,
  initialSettings,
  initialSession,
}: Props) {
  const { settings } = useSettings(initialSettings);
  const feed = settings?.metadata.feed ?? initialSettings.metadata.feed;

  switch (feed.source.type) {
    case "threads":
      return (
        <ThreadFeedScreen
          initialPage={initialData.page}
          initialPageData={initialData.threads}
          category={undefined}
          paginationBasePath="/"
          showCategorySelect={true}
          showQuickShare={feed.source.quickShare === "enabled"}
          initialSession={initialSession}
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
          initialSession={initialSession}
        />
      );
  }
}
