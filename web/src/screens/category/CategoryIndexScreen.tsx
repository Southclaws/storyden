"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import { CategoryListOKResponse } from "@/api/openapi-schema";
import {
  CategoryCardGrid,
  CategoryCardList,
} from "@/components/category/CategoryCardList/CategoryCardList";
import { useSettingsContext } from "@/components/site/SettingsContext/SettingsContext";
import { Unready } from "@/components/site/Unready";

export type Props = {
  initialCategoryList?: CategoryListOKResponse;
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
  const { feed } = useSettingsContext();
  const { ready, data, error } = useCategoryIndexScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { categories } = data;

  switch (feed.layout.type) {
    case "grid":
      return <CategoryCardGrid categories={categories} />;

    case "list":
      return <CategoryCardList categories={categories} />;
  }
}
