"use client";

import {
  CategoryList,
  NodeListResult,
  ThreadListResult,
} from "@/api/openapi-schema";
import { useSettingsContext } from "@/components/site/SettingsContext/SettingsContext";

import { CategoryIndexScreen } from "../category/CategoryIndexScreen";

import { LibraryFeedScreen } from "./LibraryFeedScreen/LibraryFeedScreen";
import { ThreadFeedScreen } from "./ThreadFeedScreen/ThreadFeedScreen";

export type InitialData = {
  threads?: ThreadListResult;
  page?: number;
  library?: NodeListResult;
  categories?: CategoryList;
};

type Props = {
  initialData: InitialData;
};

export function FeedScreenContent({ initialData }: Props) {
  const { feed } = useSettingsContext();

  switch (feed.source.type) {
    case "threads":
      return (
        <ThreadFeedScreen
          initialPage={initialData.page}
          initialPageData={initialData.threads}
          category={undefined}
          paginationBasePath="/"
        />
      );

    case "library":
      return <LibraryFeedScreen initialData={initialData.library} />;

    case "categories":
      return (
        <CategoryIndexScreen
          initialCategoryList={initialData.categories}
          initialThreadList={initialData.threads}
          initialThreadListPage={initialData.page}
          paginationBasePath="/"
        />
      );
  }
}
