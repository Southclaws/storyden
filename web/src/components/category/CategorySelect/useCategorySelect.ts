import { createListCollection } from "@ark-ui/react";

import { useCategoryList } from "@/api/openapi-client/categories";
import { ListCollectionItem } from "@/components/ui/form/SelectField";

export function useCategorySelect() {
  const { data, error } = useCategoryList();

  if (!data) {
    if (error) {
      console.error("Failed to load categories", error);
    }

    // An empty collection because we don't want to show an errored/empty state.
    const collection = createListCollection<ListCollectionItem>({ items: [] });

    return {
      ready: false as const,
      collection,
    };
  }

  const collection = createListCollection({
    items: data.categories.map((category) => ({
      label: category.name,
      value: category.id,
    })),
  });

  return {
    ready: true as const,
    collection,
  };
}
