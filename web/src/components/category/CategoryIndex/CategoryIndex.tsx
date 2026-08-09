import { Account, ThreadListResult } from "@/api/openapi-schema";
import { ComposeAnchor } from "@/components/site/Navigation/Anchors/Compose";
import { SectionHeading } from "@/components/ui/section-heading";
import { CategoryTree } from "@/lib/category/tree";
import { Settings } from "@/lib/settings/settings";
import { ThreadFeedScreen } from "@/screens/feed/ThreadFeedScreen/ThreadFeedScreen";
import { LStack, WStack } from "@/styled-system/jsx";

import { CategoryLayout } from "./CategoryCardLayout";

export type CategoryIndexContentProps = {
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

export function CategoryIndexContent({
  initialSession,
  initialSettings,
  initialThreadList,
  initialThreadListPage,
  layout,
  threadListMode,
  showQuickShare,
  categories,
  paginationBasePath,
}: CategoryIndexContentProps) {
  return (
    <LStack gap="8">
      <CategoryLayout layout={layout} categories={categories} />

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
          <SectionHeading>{heading}</SectionHeading>

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
