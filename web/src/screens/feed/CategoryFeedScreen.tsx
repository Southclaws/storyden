"use client";

import { CategoryIndexContent } from "@/components/category/CategoryIndex/CategoryIndex";
import { Unready } from "@/components/site/Unready";
import {
  type CategoryIndexScreenProps,
  useCategoryIndexScreen,
} from "@/screens/category/CategoryIndexScreen";

export function CategoryFeedScreen(props: CategoryIndexScreenProps) {
  const { ready, data, error } = useCategoryIndexScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <CategoryIndexContent
      initialThreadListPage={props.initialThreadListPage}
      initialSession={props.initialSession}
      initialSettings={props.initialSettings}
      initialThreadList={props.initialThreadList}
      layout={props.layout}
      threadListMode={props.threadListMode}
      showQuickShare={props.showQuickShare}
      categories={data.categories}
      paginationBasePath={props.paginationBasePath}
    />
  );
}
