import { createListCollection } from "@ark-ui/react";

import { useCategoryList } from "@/api/openapi-client/categories";
import { CategoryListOKResponse } from "@/api/openapi-schema";
import { SelectFieldItem } from "@/components/ui/select";

export const NO_CATEGORY_VALUE = "__none__";

export function useCategorySelect(
  initialCategoryList?: CategoryListOKResponse,
) {
  const { data, error } = useCategoryList({
    swr: {
      fallbackData: initialCategoryList,
      revalidateOnMount: initialCategoryList ? false : undefined,
    },
  });

  if (!data) {
    if (error) {
      console.error("Failed to load categories", error);
    }

    // An empty collection because we don't want to show an errored/empty state.
    const collection = createListCollection<SelectFieldItem>({ items: [] });

    return {
      ready: false as const,
      collection,
      error,
    };
  }

  const collection = createListCollection({
    items: [
      { label: "No category", value: NO_CATEGORY_VALUE },
      ...data.categories.map((category) => ({
        label: category.name,
        value: category.id,
      })),
    ],
  });

  return {
    ready: true as const,
    collection,
  };
}
