"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import {
  Account,
  CategoryListOKResponse,
  ThreadListResult,
} from "@/api/openapi-schema";
import { CategoryIndex } from "@/components/category/CategoryIndex/CategoryIndex";
import { Unready } from "@/components/site/Unready";
import { buildCategoryTree } from "@/lib/category/tree";
import { Settings } from "@/lib/settings/settings";

export type Props = {
  initialThreadListPage?: number;
  initialThreadList?: ThreadListResult;
  initialSession?: Account;
  initialSettings?: Settings;
  initialCategoryList?: CategoryListOKResponse;

  layout: "grid" | "list";
  threadListMode: "none" | "all" | "uncategorised";
  showQuickShare: boolean;
  paginationBasePath: string;
};

export function useCategoryIndexScreen({ initialCategoryList }: Props) {
  const { data, error } = useCategoryList({
    swr: { fallbackData: initialCategoryList },
  });
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data: {
      categories: data.categories,
    },
  };
}

export function CategoryIndexScreen(props: Props) {
  const { ready, data, error } = useCategoryIndexScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { categories } = data;

  const tree = buildCategoryTree(categories);

  return (
    <CategoryIndex
      initialThreadListPage={props.initialThreadListPage}
      initialSession={props.initialSession}
      initialSettings={props.initialSettings}
      initialThreadList={props.initialThreadList}
      layout={props.layout}
      threadListMode={props.threadListMode}
      showQuickShare={props.showQuickShare}
      categories={tree}
      paginationBasePath={props.paginationBasePath}
    />
  );
}
