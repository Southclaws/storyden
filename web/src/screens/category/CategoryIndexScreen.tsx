"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import {
  Account,
  CategoryListOKResponse,
  ThreadListResult,
} from "@/api/openapi-schema";
import { CategoryBadge } from "@/components/category/CategoryBadge";
import { CategoryCreateTrigger } from "@/components/category/CategoryCreate/CategoryCreateTrigger";
import { CategoryIndexContent } from "@/components/category/CategoryIndex/CategoryIndex";
import { Unready } from "@/components/site/Unready";
import { PageHeader } from "@/components/ui/page-header";
import { buildCategoryTree } from "@/lib/category/tree";
import { Settings } from "@/lib/settings/settings";
import { HStack, LStack } from "@/styled-system/jsx";

export type CategoryIndexScreenProps = {
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

export function useCategoryIndexScreen({
  initialCategoryList,
}: CategoryIndexScreenProps) {
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
      categories: buildCategoryTree(data.categories),
    },
  };
}

export function CategoryIndexScreen(props: CategoryIndexScreenProps) {
  const { ready, data, error } = useCategoryIndexScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { categories } = data;

  const categoryCount = categories.length;

  return (
    <LStack gap="4">
      <LStack gap="3">
        <PageHeader
          title="Discussion categories"
          description={getCategoryDescription(categoryCount)}
          actions={<CategoryCreateTrigger />}
        />

        <HStack gap="2" flexWrap="wrap">
          {categories.map((category) => (
            <CategoryBadge key={category.id} category={category} />
          ))}
        </HStack>
      </LStack>

      <CategoryIndexContent
        initialThreadListPage={props.initialThreadListPage}
        initialSession={props.initialSession}
        initialSettings={props.initialSettings}
        initialThreadList={props.initialThreadList}
        layout={props.layout}
        threadListMode={props.threadListMode}
        showQuickShare={props.showQuickShare}
        categories={categories}
        paginationBasePath={props.paginationBasePath}
      />
    </LStack>
  );
}

function getCategoryDescription(categoryCount: number) {
  if (categoryCount === 0) {
    return "No categories yet. Create one?";
  }

  if (categoryCount === 1) {
    return "There is 1 category available to start a discussion.";
  }

  return `There are ${categoryCount} categories available to start discussions.`;
}
