import { FeedBlock } from "@/lib/settings/feed";
import { ThreadFeed } from "@/screens/feed/ThreadFeedScreen/ThreadFeedScreen";

import { useFeedBlockEditor } from "../Context";

export function FeedThreadsBlock({
  block,
}: {
  block: Extract<FeedBlock, { type: "threads" }>;
}) {
  const { initialData } = useFeedBlockEditor();

  return (
    <ThreadFeed
      initialPage={initialData.initialPage}
      initialPageData={
        initialData.initialThreadSource === block.source
          ? initialData.initialThreadList
          : undefined
      }
      category={block.source === "all" ? undefined : null}
      paginationBasePath="/"
    />
  );
}
