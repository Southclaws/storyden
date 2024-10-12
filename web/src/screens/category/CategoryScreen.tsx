"use client";

import { useCategoryList } from "@/api/openapi-client/categories";
import { useThreadList } from "@/api/openapi-client/threads";
import {
  CategoryListOKResponse,
  Permission,
  ThreadListOKResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { CategoryMenu } from "@/components/category/CategoryMenu/CategoryMenu";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { HStack, LStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { FeedScreen } from "../feed/FeedScreen";

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

export function CategoryScreen(props: Props) {
  const { ready, data, error } = useCategoryScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { category, threads } = data;

  return (
    <LStack>
      <LStack gap="1">
        <HStack w="full" justify="space-between" alignItems="start">
          <Heading>{category.name}</Heading>

          <CategoryMenu category={category} />
        </HStack>

        <styled.p color="fg.muted">{category.description}</styled.p>
      </LStack>

      <FeedScreen
        params={{
          categories: [props.slug],
        }}
        initialData={threads}
      />
    </LStack>
  );
}
