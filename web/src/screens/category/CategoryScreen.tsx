"use client";

import { useCategoryGet } from "@/api/openapi-client/categories";
import {
  CategoryGetOKResponse,
  Permission,
  ThreadListOKResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { CategoryLayout } from "@/components/category/CategoryIndex/CategoryCardLayout";
import { CategoryMenu } from "@/components/category/CategoryMenu/CategoryMenu";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { Box, LStack, WStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { hasPermission } from "@/utils/permissions";

import { ThreadFeedScreen } from "../feed/ThreadFeedScreen/ThreadFeedScreen";

export type Props = {
  initialCategory: CategoryGetOKResponse;
  initialThreadList: ThreadListOKResponse;
  slug: string;
};

export function useCategoryScreen({ initialCategory, slug }: Props) {
  const session = useSession();

  const { data, error } = useCategoryGet(slug, {
    swr: { fallbackData: initialCategory },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const canEditCategory = hasPermission(session, Permission.MANAGE_CATEGORIES);

  return {
    ready: true as const,
    data: {
      canEditCategory,
      category: data,
    },
  };
}

type ScreenProps = {
  initialPage: number;
} & Props;

export function CategoryScreen(props: ScreenProps) {
  const { ready, data, error } = useCategoryScreen(props);
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const { category } = data;
  const coverImageURL = getAssetURL(category.cover_image?.path);

  return (
    <LStack>
      {coverImageURL && (
        <Box height="auto" width="full">
          <styled.img
            src={coverImageURL}
            alt="" // No alt image, decorative
            aria-hidden="true"
            width="full"
            height="full"
            borderRadius="md"
            objectFit="cover"
            objectPosition="center"
          />
        </Box>
      )}

      <LStack gap="1">
        <WStack alignItems="start">
          <Heading>{category.name}</Heading>

          <CategoryMenu category={category} />
        </WStack>

        <styled.p color="fg.muted">{category.description}</styled.p>
      </LStack>

      {category.children && category.children.length > 0 && (
        <LStack gap="1">
          <Heading size="sm" color="fg.muted">
            Subcategories
          </Heading>
          <CategoryLayout layout="grid" categories={category.children} />
        </LStack>
      )}

      <ThreadFeedScreen
        initialPage={props.initialPage}
        initialPageData={props.initialThreadList}
        category={category}
        paginationBasePath={`/d/${data.category.slug}`}
        showCategorySelect={false}
        hideCategoryBadge={true}
      />
    </LStack>
  );
}
