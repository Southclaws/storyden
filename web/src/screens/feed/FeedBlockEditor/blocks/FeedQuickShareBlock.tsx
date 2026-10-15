import { QuickShare } from "@/components/feed/QuickShare/QuickShare";
import { FeedBlock } from "@/lib/settings/feed";

import { useFeedBlockEditor } from "../Context";

export function FeedQuickShareBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "quick-share" }>;
}) {
  const { initialSession, initialData } = useFeedBlockEditor();

  return (
    <QuickShare
      initialSession={initialSession}
      initialCategory={undefined}
      initialCategoryList={initialData.initialCategoryList}
      showCategorySelect={block.showCategorySelect}
    />
  );
}
