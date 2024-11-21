"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import { useThreadList } from "@/api/openapi-client/threads";
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

export function useCategoryScreen({
  initialCategoryList,
  initialThreadList,
  slug,
}: Props) {
  const session = useSession();

  const { data: categoryListData, error: categoryListError } = useCategoryList({
    swr: { fallbackData: initialCategoryList },
  });

  const { data: threadListData, error: threadListError } = useThreadList(
    { categories: [slug] },
    {
      swr: { fallbackData: initialThreadList },
    },
  );

  if (!categoryListData) {
    return {
      ready: false as const,
      error: categoryListError,
    };
  }
  if (!threadListData) {
    return {
      ready: false as const,
      error: threadListError,
    };
  }

  const category = categoryListData.categories.find((c) => c.slug === slug);
  if (!category) {
    return {
      ready: false as const,
      error: new Error("Category not found"),
    };
  }

  const canEditCategory = hasPermission(session, Permission.MANAGE_CATEGORIES);

  const threads = threadListData;

  return {
    ready: true as const,
    data: {
      canEditCategory,
      category,
      threads,
    },
  };
}
