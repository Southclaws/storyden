import { match } from "ts-pattern";

import { ThreadListResult } from "@/api/openapi-schema";
import { ComposeAnchor } from "@/components/site/Navigation/Anchors/Compose";
import { Heading } from "@/components/ui/heading";
import { CategoryTree } from "@/lib/category/tree";
import { ThreadFeedScreen } from "@/screens/feed/ThreadFeedScreen/ThreadFeedScreen";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { CategoryBadge } from "../CategoryBadge";
import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";

import { CategoryLayout } from "./CategoryCardLayout";

export type Props = {
  layout: "grid" | "list";
  threadListMode: "none" | "all" | "uncategorised";
  showQuickShare: boolean;
  categories: CategoryTree[];
  initialThreadList?: ThreadListResult;
  initialThreadListPage?: number;
  paginationBasePath: string;
};

export function CategoryIndex({
  layout,
  threadListMode,
  showQuickShare,
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
        mode={threadListMode}
        showQuickShare={showQuickShare}
        initialThreadList={initialThreadList}
        initialPage={initialThreadListPage}
        paginationBasePath={paginationBasePath}
      />
    </LStack>
  );
}

function ThreadListSection({
  mode,
  showQuickShare,
  initialThreadList,
  initialPage,
  paginationBasePath,
}: {
  mode: "none" | "all" | "uncategorised";
  showQuickShare: boolean;
  initialThreadList?: ThreadListResult;
  initialPage?: number;
  paginationBasePath: string;
}) {
  if (mode === "none") {
    return null;
  }

  const heading =
    mode === "all"
      ? "All discussion threads"
      : "Uncategorised discussion threads";

  // Only show the category select when showing all threads, not uncategorised.
  const showCategorySelect = mode === "all";

  return (
    <LStack>
      {!showQuickShare && (
        <WStack>
          <Heading>{heading}</Heading>

          <ComposeAnchor />
        </WStack>
      )}

      <ThreadFeedScreen
        initialPage={initialPage}
        initialPageData={initialThreadList}
        category={mode === "all" ? undefined : null}
        paginationBasePath={paginationBasePath}
        showCategorySelect={showCategorySelect}
        showQuickShare={showQuickShare}
      />
    </LStack>
  );
}
