import Link from "next/link";
import { match } from "ts-pattern";

import { ThreadListResult } from "@/api/openapi-schema";
import { ComposeAnchor } from "@/components/site/Navigation/Anchors/Compose";
import { Heading } from "@/components/ui/heading";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { categoryColourCSS } from "@/lib/category/colours";
import { CategoryTree } from "@/lib/category/tree";
import { ThreadFeedScreen } from "@/screens/feed/ThreadFeedScreen/ThreadFeedScreen";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";

import { CategoryBadge } from "../CategoryBadge";
import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";
import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

import { CategoryCardGrid, CategoryLayout } from "./CategoryCardLayout";

export type Props = {
  layout: "grid" | "list";
  categories: CategoryTree[];
  initialThreadList?: ThreadListResult;
  initialThreadListPage?: number;
  paginationBasePath: string;
};

export function CategoryIndex({
  layout,
  categories,
  initialThreadList,
  initialThreadListPage,
  paginationBasePath,
}: Props) {
  const categoryCount = categories.length;

  return (
    <LStack gap="8">
      <LStack>
        <WStack>
          <Heading>Discussion categories</Heading>

          <CategoryCreateTrigger />
        </WStack>

        <LStack>
          {match(categoryCount)
            .when(
              (c) => c === 0,
              () => (
                <styled.p color="fg.muted">
                  No categories yet. Create one?
                </styled.p>
              ),
            )
            .when(
              (c) => c === 1,
              () => (
                <styled.p color="fg.muted">
                  There is {categoryCount} category available to start a
                  discussion.
                </styled.p>
              ),
            )
            .otherwise(() => (
              <styled.p color="fg.muted">
                There are {categoryCount} categories available to start
                discussions.
              </styled.p>
            ))}

          <HStack>
            {categories.map((c) => (
              <CategoryBadge key={c.id} category={c} />
            ))}
          </HStack>
        </LStack>

        <CategoryLayout layout={layout} categories={categories} />
      </LStack>

      <ThreadListSection
        initialThreadList={initialThreadList}
        initialPage={initialThreadListPage}
        paginationBasePath={paginationBasePath}
      />
    </LStack>
  );
}

function ThreadListSection({
  initialThreadList,
  initialPage,
  paginationBasePath,
}: {
  initialThreadList?: ThreadListResult;
  initialPage?: number;
  paginationBasePath: string;
}) {
  if (!initialThreadList?.threads) {
    return null;
  }

  return (
    <LStack>
      <WStack>
        <Heading>Uncategorised discussion threads</Heading>

        <ComposeAnchor />
      </WStack>

      <styled.p color="fg.muted">
        Threads that have not been posted within a category.
      </styled.p>

      <ThreadFeedScreen
        initialPage={initialPage}
        initialPageData={initialThreadList}
        category={null}
        paginationBasePath={paginationBasePath}
        showCategorySelect={false}
      />
    </LStack>
  );
}
