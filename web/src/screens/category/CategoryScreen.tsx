"use client";

import { CategoryMenu } from "@/components/category/CategoryMenu/CategoryMenu";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { LStack, WStack, styled } from "@/styled-system/jsx";

import { ThreadFeedScreen } from "../feed/ThreadFeedScreen/ThreadFeedScreen";

import { Props, useCategoryScreen } from "./useCategoryScreen";

type ScreenProps = {
  initialPage: number;
} & Props;

export function CategoryScreen(props: ScreenProps) {
  const { ready, data, error } = useCategoryScreen(props);
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  const { category } = data;

  return (
    <LStack>
      <LStack gap="1">
        <WStack alignItems="start">
          <Heading>{category.name}</Heading>

          <CategoryMenu category={category} />
        </WStack>

        <styled.p color="fg.muted">{category.description}</styled.p>
      </LStack>

      <ThreadFeedScreen
        initialPage={props.initialPage}
        initialPageData={props.initialThreadList}
        category={category}
        paginationBasePath={`/d/${data.category.slug}`}
        showCategorySelect={false}
      />
    </LStack>
  );
}
