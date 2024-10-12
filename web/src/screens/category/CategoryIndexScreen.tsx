"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import { CategoryListOKResponse } from "@/api/openapi-schema";
import { CategoryCardList } from "@/components/category/CategoryCardList/CategoryCardList";
import { Unready } from "@/components/site/Unready";

export type Props = {
  initialCategoryList: CategoryListOKResponse;
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

  return <CategoryCardList categories={categories} />;
}
