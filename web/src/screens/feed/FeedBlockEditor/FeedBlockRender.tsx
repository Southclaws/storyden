import { ReactNode } from "react";

import { FeedBlock } from "@/lib/settings/feed";
import { Box } from "@/styled-system/jsx";

import { FeedCategoriesBlock } from "./blocks/FeedCategoriesBlock";
import { FeedContentBlock } from "./blocks/FeedContentBlock";
import { FeedCoverBlock } from "./blocks/FeedCoverBlock";
import { FeedLibraryBlock } from "./blocks/FeedLibraryBlock";
import { FeedQuickShareBlock } from "./blocks/FeedQuickShareBlock";
import { FeedSubtitleBlock } from "./blocks/FeedSubtitleBlock";
import { FeedThreadsBlock } from "./blocks/FeedThreadsBlock";
import { FeedTitleBlock } from "./blocks/FeedTitleBlock";

export function FeedBlockRender({
  block,
  children,
}: {
  block: FeedBlock;
  children?: ReactNode;
}) {
  return (
    <Box
      className={"feed-page__block feed-page__block--" + block.type}
      data-sd-block={block.type}
      width="full"
    >
      {children ?? <FeedBlockViewContent block={block} />}
    </Box>
  );
}

function FeedBlockViewContent({ block }: { block: FeedBlock }) {
  switch (block.type) {
    case "title":
      return <FeedTitleBlock />;
    case "subtitle":
      return <FeedSubtitleBlock />;
    case "content":
      return <FeedContentBlock />;
    case "cover":
      return <FeedCoverBlock />;
    case "categories":
      return <FeedCategoriesBlock block={block} />;
    case "threads":
      return <FeedThreadsBlock block={block} />;
    case "quick-share":
      return <FeedQuickShareBlock block={block} />;
    case "library":
      return <FeedLibraryBlock block={block} />;
  }
}
