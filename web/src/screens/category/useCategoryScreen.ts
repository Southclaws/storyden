"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import {
  CategoryListOKResponse,
  Permission,
  ThreadListOKResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  initialCategoryList: CategoryListOKResponse;
  initialThreadList: ThreadListOKResponse;
  slug: string;
};

export function useCategoryScreen({ initialCategoryList, slug }: Props) {
  const session = useSession();

  const { data, error } = useCategoryList({
    swr: { fallbackData: initialCategoryList },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  // TODO: Write a get category endpoint
  const category = data.categories.find((c) => c.slug === slug);
  if (!category) {
    return {
      ready: false as const,
      error: new Error("Category not found"),
    };
  }

  const canEditCategory = hasPermission(session, Permission.MANAGE_CATEGORIES);

  return {
    ready: true as const,
    data: {
      canEditCategory,
      category,
    },
  };
}
