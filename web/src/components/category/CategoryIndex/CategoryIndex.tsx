import { match } from "ts-pattern";

import { Account, ThreadListResult } from "@/api/openapi-schema";
import { ComposeAnchor } from "@/components/site/Navigation/Anchors/Compose";
import { Heading } from "@/components/ui/heading";
import { CategoryTree } from "@/lib/category/tree";
import { useI18n } from "@/i18n/provider";
import { Settings } from "@/lib/settings/settings";
import { ThreadFeedScreen } from "@/screens/feed/ThreadFeedScreen/ThreadFeedScreen";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { CategoryBadge } from "../CategoryBadge";
import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";

import { CategoryLayout } from "./CategoryCardLayout";

export type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
  initialThreadList?: ThreadListResult;
  initialThreadListPage?: number;

  layout: "grid" | "list";
  threadListMode: "none" | "all" | "uncategorised";
  showQuickShare: boolean;
  categories: CategoryTree[];
  paginationBasePath: string;
};

export function CategoryIndex({
  initialSession,
  initialSettings,
  initialThreadList,
  initialThreadListPage,
  layout,
  threadListMode,
  showQuickShare,
  categories,
  paginationBasePath,
}: Props) {
  const { t } = useI18n();
  const categoryCount = categories.length;

  return (
    <LStack gap="8">
      <LStack>
        <WStack>
          <Heading>{t("Discussion categories")}</Heading>

          <CategoryCreateTrigger />
        </WStack>

        <LStack>
          {match(categoryCount)
            .when(
              (c) => c === 0,
              () => (
                <styled.p color="fg.muted">
                  {t("No categories yet. Create one?")}
                </styled.p>
              ),
            )
            .when(
              (c) => c === 1,
              () => (
                <styled.p color="fg.muted">
                  {categoryCount}{" "}
                  {t("category available to start a discussion.")}
                </styled.p>
              ),
            )
            .otherwise(() => (
              <styled.p color="fg.muted">
                {categoryCount}{" "}
                {t("categories available to start discussions.")}
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
        initialThreadListPage={initialThreadListPage}
        initialSession={initialSession}
        initialSettings={initialSettings}
        initialThreadList={initialThreadList}
        mode={threadListMode}
        showQuickShare={showQuickShare}
        paginationBasePath={paginationBasePath}
      />
    </LStack>
  );
}

function ThreadListSection({
  initialSession,
  initialSettings,
  initialThreadList,
  initialThreadListPage,
  mode,
  showQuickShare,
  paginationBasePath,
}: {
  initialSession?: Account;
  initialSettings?: Settings;
  initialThreadList?: ThreadListResult;
  initialThreadListPage?: number;
  mode: "none" | "all" | "uncategorised";
  showQuickShare: boolean;
  paginationBasePath: string;
}) {
  const { t } = useI18n();
  if (mode === "none") {
    return null;
  }

  const heading =
    mode === "all"
      ? t("All discussion threads")
      : t("Uncategorised discussion threads");

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
        initialPage={initialThreadListPage}
        initialPageData={initialThreadList}
        initialSession={initialSession}
        initialSettings={initialSettings}
        category={mode === "all" ? undefined : null}
        paginationBasePath={paginationBasePath}
        showCategorySelect={showCategorySelect}
        showQuickShare={showQuickShare}
      />
    </LStack>
  );
}
