import { QuickShare } from "@/components/feed/QuickShare/QuickShare";
import { FeedBlock } from "@/lib/settings/feed";

import { useFeedBlockEditor } from "../Context";

export function FeedQuickShareBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "quick-share" }>;
}) {
  const { initialSession } = useFeedBlockEditor();

  return (
    <QuickShare
      initialSession={initialSession}
      initialCategory={undefined}
      showCategorySelect={block.showCategorySelect}
    />
  );
}
