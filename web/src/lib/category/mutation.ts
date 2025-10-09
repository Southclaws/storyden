import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import {
  categoryDelete,
  categoryUpdate,
  getCategoryListKey,
} from "@/api/openapi-client/categories";
import {
  CategoryDeleteBody,
  CategoryListOKResponse,
  CategoryMutableProps,
} from "@/api/openapi-schema";
import { categoryListResponse } from "@/api/openapi-server/categories";

function threadsKeyFilter(key: Arguments) {
  if (!Array.isArray(key)) return false;

  if (key.length === 0) {
    return false;
  }

  const path = key[0];
  if (typeof path !== "string") {
    return false;
  }

  if (!path.startsWith("/threads")) {
    return false;
  }

  return true;
}

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

  const updateCategory = async (
    slug: string,
    updated: CategoryMutableProps,
  ) => {
    const mutator: MutatorCallback<CategoryListOKResponse> = (data) => {
      if (!data) return;

      const newData = {
        categories: data.categories.map((category) => {
          if (category.slug === slug) {
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

    await categoryUpdate(slug, updated);
  };

  const deleteCategory = async (slug: string, body: CategoryDeleteBody) => {
    await categoryDelete(slug, body);
    await mutate(keyFilterFn);
    await mutate(threadsKeyFilter);
  };

  return {
    revalidateList,
    updateCategory,
    deleteCategory,
  };
}
