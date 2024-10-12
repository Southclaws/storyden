import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import {
  categoryUpdate,
  getCategoryListKey,
} from "@/api/openapi-client/categories";
import {
  CategoryListOKResponse,
  CategoryMutableProps,
} from "@/api/openapi-schema";
import { categoryListResponse } from "@/api/openapi-server/categories";

export function useCategoryMutations() {
  const { mutate } = useSWRConfig();

  const categoryGroupKey = getCategoryListKey();

  function keyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(categoryGroupKey);
  }

  const revalidateList = async (
    data?: MutatorCallback<categoryListResponse>,
  ) => {
    await mutate(keyFilterFn, data);
  };

  const updateCategory = async (id: string, updated: CategoryMutableProps) => {
    const mutator: MutatorCallback<CategoryListOKResponse> = (data) => {
      if (!data) return;

      const newData = {
        categories: data.categories.map((category) => {
          if (category.id === id) {
            return {
              ...category,
              ...updated,
            };
          }

          return category;
        }),
      };

      return newData;
    };

    await mutate(categoryGroupKey, mutator, {
      revalidate: false,
    });

    await categoryUpdate(id, updated);
  };

  return {
    revalidateList,
    updateCategory,
  };
}
